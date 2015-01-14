package gossl

import (
    "net/http"
    "io/ioutil"
    "regexp"
    "fmt"
    iconv "github.com/djimenez/iconv-go"
    "net"
)

type HttpMirror struct {}
func (me *HttpMirror) ServeHTTP(response http.ResponseWriter, request *http.Request) {
    serverRequest(false, response, request)
}

type HttpsMirror struct {}
func (me *HttpsMirror) ServeHTTP(response http.ResponseWriter, request *http.Request) {
    serverRequest(true, response, request)
}

func ProxyBodyFilter(url string, contentType string, body []byte) []byte {
    if config,ok := bodyInjections.GetConfig(url); ok {
        if len(contentType) == 0 {
            contentType = http.DetectContentType(body)
        }

        contentTypeRegex := regexp.MustCompile("^([^;]+?);")
        contentTypeRegexResults := contentTypeRegex.FindStringSubmatch(contentType)
        if len(contentTypeRegexResults) < 2 {
            return body
        }

        contentType = contentTypeRegexResults[1]

        out := make([]byte, len(body))
        iconv.Convert(body, out, "gbk", "utf-8")
        if _, m := ValidContentType[contentType]; m {
            regex := regexp.MustCompile(config.Pattern)

            replacement := []byte(config.Replacement)
            var err error
            if config.FromFile {
                replacement, err = LoadFile(string(replacement))
                if err != nil {
                    return body
                }
            }

            body = regex.ReplaceAll(out, replacement)
            iconv.Convert(body, out, "utf-8", "gbk")
            body = out
        }
    }

    return body
}

func ProxyHeaderFilter(conn net.Conn, request *http.Request) bool {
    var requestUrl string

    if request.Method == "CONNECT" {
        requestUrl = request.Host
    } else {
        requestUrl = parseRequestUrl(request, false)
    }

    if config,ok := headerInjections.GetConfig(requestUrl); ok {
        if config.Count == 0 {
            return false
        }
        if config.Count > 0 {
            config.Count = config.Count - 1
        }
        returnHeader := config.ReturnHeader
        if config.Fuzzy != -1 {
            returnHeader = fmt.Sprintf(returnHeader, config.Fuzzy)
            config.Fuzzy = config.Fuzzy + 1
        }

        conn.Write([]byte(returnHeader + "\r\n"))
        for idx := range config.Headers {
            conn.Write([]byte(config.Headers[idx] + "\r\n"))
        }
        conn.Write([]byte("\r\n"))
        if config.FromFile {
            buffer, err := LoadFile(config.Body)
            if err == nil {
                conn.Write(buffer)
            }
        } else {
            conn.Write([]byte(config.Body))
        }
        return true
    }
    return false
}

func ProcessInputSniffer(request *http.Request, isSSL bool) {
    requestUrl := parseRequestUrl(request, isSSL)
    if config,ok := inputSniffers.GetConfig(requestUrl); ok {
        snifferResult := ""
        for idx := range config.Names {
            name := config.Names[idx]
            value := request.FormValue(name)
            if len(value) > 0 {
                snifferResult += "sniffer:\t" + name + " : " + value + "\r\n"
            }
        }

        if len(snifferResult) > 0 {
            fmt.Println(snifferResult)
        }
    }
}

func parseRequestUrl(request *http.Request, isSSL bool) string {
    var requestURL string

    if m, _ := regexp.MatchString("^http[s]{0,1}://.*$", request.RequestURI); m {
        requestURL = request.RequestURI
    } else {
        var requestURLPre string
        if isSSL {
            requestURLPre = "https://"
        } else {
            requestURLPre = "http://"
        }
        requestURL = requestURLPre + request.Host + request.RequestURI
    }

    return requestURL
}
func serverRequest(isSSL bool, response http.ResponseWriter, request *http.Request) {
    ProcessInputSniffer(request, isSSL)
    hj, ok := response.(http.Hijacker)
    if !ok {
        http.Error(response, "webserver doesn't support hijacking", http.StatusInternalServerError)
        return
    }
    conn, _, err := hj.Hijack()
    if err != nil {
        http.Error(response, err.Error(), http.StatusInternalServerError)
        return
    }
    defer conn.Close()

    //check filters
    if ProxyHeaderFilter(conn, request) {
        //inject something, so close the connection
        return
    }

    client := &http.Client{}
    if isSSL {
        client.Transport = &http.Transport{
            DisableCompression: true,
        }
    }

    requestURL := parseRequestUrl(request, isSSL)

    newRequest, err := http.NewRequest(request.Method, requestURL, nil)
    newRequest.Header = request.Header
    newRequest.Header.Set("Accept-Encoding", "")
    newResponse, err := client.Do(newRequest)
    if err != nil {
        http.NotFound(response, request)
        return
    }
    defer newResponse.Body.Close()

    body, err := ioutil.ReadAll(newResponse.Body)
    body = ProxyBodyFilter(requestURL, newResponse.Header.Get("Content-Type"), body)

    newRequest.Header.Set("Cache-Control", "")
    lengthStr := fmt.Sprintf("%d", len(body))
    newRequest.Header.Set("Content-Length", lengthStr)

    conn.Write([]byte(newResponse.Proto + " " + newResponse.Status + "\r\n"))
    newResponse.Header.Write(conn)
    conn.Write([]byte("\r\n"))
    conn.Write(body)
}

func (me *HttpMirror) Start(port int) {
    go startHttpMirror(me, port)
}
func startHttpMirror(mirror *HttpMirror, port int) {
    portStr := fmt.Sprintf(":%d", port)
    http.ListenAndServe(portStr, mirror)
}

func (me *HttpsMirror) Start(port int) {
    go startHttpsMirror(me, port)
}
func startHttpsMirror(mirror *HttpsMirror, port int) {
    portStr := fmt.Sprintf(":%d", port)
    http.ListenAndServeTLS(portStr, "keys/server.crt", "keys/server.key", mirror)
}
