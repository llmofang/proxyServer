package gossl

import (
    "os"
    "encoding/json"
    "errors"
    "fmt"
    "regexp"
)

type ConfigsInterface interface{
    GetConfig()
}

//sniffer
type InputSniffer struct {
    Url string
    UrlRegex bool
	Names []string
    Count int
}
type InputSnifferConfigs struct {
    Configs []InputSniffer
}
func (me *InputSnifferConfigs) GetConfig(url string) (*InputSniffer, bool) {
    for idx := range me.Configs {
        tmpConfig := &me.Configs[idx]
        if tmpConfig.UrlRegex {
            if m, _ := regexp.MatchString(tmpConfig.Url, url); m {
                return tmpConfig, true
            }
        } else {
            if tmpConfig.Url == url {
                return tmpConfig, true
            }
        }
    }
    return nil, false
}
var inputSniffers InputSnifferConfigs

//header injections
type HeaderInjection struct {
    Url string
    UrlRegex bool
	ReturnHeader string
    Fuzzy int
    Headers []string
    FromFile bool
    Body string
    Count int
}
type HeaderInjectionConfigs struct {
    Configs []HeaderInjection
}
func (me *HeaderInjectionConfigs) GetConfig(url string) (*HeaderInjection, bool) {
    for idx := range me.Configs {
        tmpConfig := &me.Configs[idx]
        if tmpConfig.UrlRegex {
            if m, _ := regexp.MatchString(tmpConfig.Url, url); m {
                return tmpConfig, true
            }
        } else {
            if tmpConfig.Url == url {
                return tmpConfig, true
            }
        }
    }
    return nil, false
}
var headerInjections HeaderInjectionConfigs

//header injections
type BodyInjection struct {
    Url string
    UrlRegex bool
    Pattern string
    FromFile bool
	Replacement string
    Flag string
    Count int
}
type BodyInjectionConfigs struct {
    Configs []BodyInjection
}
func (me *BodyInjectionConfigs) GetConfig(url string) (*BodyInjection, bool) {
    for idx := range me.Configs {
        tmpConfig := &me.Configs[idx]
        if tmpConfig.UrlRegex {
            if m, _ := regexp.MatchString(tmpConfig.Url, url); m {
                return tmpConfig, true
            }
        } else {
            if tmpConfig.Url == url {
                return tmpConfig, true
            }
        }
    }
    return nil, false
}
var bodyInjections BodyInjectionConfigs


var ValidContentType = map[string]bool{"text/html": true,
                                        "text/plain": true,
                                        "application/javascript": true,
                                        "text/json": true,
                                        "text/xml": true}

func LoadFile(filePath string) ([]byte, error) {
    file, err := os.Open(filePath) // For read access.
    if err != nil {
        fmt.Println(err)
        return nil, errors.New("read file error")
    }
    defer file.Close()

    fileInfo, err := file.Stat()
    if err != nil {
        fmt.Println(err)
        return nil, errors.New("read file error")
    }

    buffer := make([]byte, fileInfo.Size())
    file.Read(buffer)

    return buffer, nil
}
func LoadJSON(filePath string, configObj interface{}) (error) {
    buffer, err := LoadFile(filePath)
    if err != nil {
        fmt.Println(err)
        return errors.New("read file error")
    }

    err = json.Unmarshal(buffer, &configObj)
    if err != nil {
		fmt.Println("error:", err)
        return errors.New("convert config json error")
	}

    return nil
}

func loadHeaderInjections(prefix string) {
    LoadJSON(prefix + "/headerinjection.js", &headerInjections.Configs)
}

func loadBodyInjections(prefix string) {
    LoadJSON(prefix + "/bodyinjection.js", &bodyInjections.Configs)
}

func loadInputSniffers(prefix string) {
    LoadJSON(prefix + "/inputsniffer.js", &inputSniffers.Configs)
}

func LoadConfigs(prefix string) {
    loadHeaderInjections(prefix)
    loadBodyInjections(prefix)
    loadInputSniffers(prefix)
}

