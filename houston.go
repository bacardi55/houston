package main

import (
  "fmt"
  "time"
  "os"
  "context"
  "strings"

  "net/url"

  "git.sr.ht/~adnano/go-gemini"
)

func main() {
    if os.Getenv("QUERY_STRING") == "" {
        fmt.Println("10\ttext/gemtext\r\n")
        fmt.Println("Enter the URL to test: ")
        os.Exit(0)
    } else {
        fmt.Println("20\ttext/gemini\r\n")
    }

    remoteUrl := os.Getenv("QUERY_STRING")

    u, e := validateUrl(remoteUrl);
    if e != nil {
        fmt.Printf("Url is not valid.\r\n")
        endResponse()
    }

    req, err := fetchGeminiPage(u)

    if err != nil {
        fmt.Printf("The url you are testing (%v) seems down.\r\n", u)
        endResponse()
    }

    // if the server response if an error code, capsule isn't up.
    if req.Status == gemini.StatusTemporaryFailure || req.Status == gemini.StatusServerUnavailable || req.Status == gemini.StatusCGIError || req.Status == gemini.StatusProxyError || req.Status == gemini.StatusPermanentFailure || req.Status == gemini.StatusGone || req.Status == gemini.StatusProxyRequestRefused || req.Status == gemini.StatusBadRequest {
        fmt.Printf("The url you are testing (%v) seems down, status is %v\r\n", u, req.Status)
        endResponse()
    }

    // Todo: Should I follow redirect and check redirect status' code…?
    // Or is a redirect response enough to see that the capsule is up?

    fmt.Printf("Everything seems ok with %v\r\n", u)
    endResponse()
}

func fetchGeminiPage(remoteUrl string) (*gemini.Response, error) {
    // Todo: configurable:
    maxRequestTime := 10

    gemclient := &gemini.Client{}
    ctx, _ := context.WithTimeout(context.Background(), time.Duration(maxRequestTime)*time.Second)
    response, err := gemclient.Get(ctx, remoteUrl)

    if err != nil {
        return nil, err
    }

    return response, nil
}

func validateUrl(remoteUrl string) (string, error) {
    remote, e := url.QueryUnescape(remoteUrl)
    if e != nil {
        return "", fmt.Errorf("Provided URL is not a good URL: %s", e)
    }
    remote = strings.Replace(remote, "..", "", -1)

    u, err := url.Parse(remote)

    if err != nil {
        return "", fmt.Errorf("Provided URL is not a good URL: %s", err)
    } else if u.Scheme != "gemini" && u.Scheme != "" {
        return "", fmt.Errorf("Only gemini url are supported for now.")
    } else {
        return "gemini://" + u.Host + u.Path, nil
    }

    fmt.Println(remote)
    return remote, nil
}

func endResponse() {
    fmt.Printf("=> /index.gmi Return to Houston homepage\r\n")
    fmt.Printf("=> /cgi-bin/houston Test another URL\r\n")
    os.Exit(0)
}

