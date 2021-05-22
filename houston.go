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
        fmt.Println("20\ttext/gemtext\r\n")
    }

    remoteUrl := os.Getenv("QUERY_STRING")

    u, e := validateUrl(remoteUrl);
    if e != nil {
        fmt.Printf("Url is not valid:\n %v\r\n", e)
        os.Exit(0)
    }

    req, err := fetchGeminiPage(u)

    if err != nil {
        fmt.Printf("The url you are testing (%v) seems down:\n%v\r\n", u, err)
        os.Exit(0)
    }

    // if the server response if an error code, capsule isn't up.
    if req.Status == gemini.StatusTemporaryFailure || req.Status == gemini.StatusServerUnavailable || req.Status == gemini.StatusCGIError || req.Status == gemini.StatusProxyError || req.Status == gemini.StatusPermanentFailure || req.Status == gemini.StatusGone || req.Status == gemini.StatusProxyRequestRefused || req.Status == gemini.StatusBadRequest {
        fmt.Printf("The url you are testing (%v) seems down:\n%v\r\n", u, err)
    }

    // Todo: Should I follow redirect and check redirect status' code…?
    // Or is a redirect response enough to see that the capsule is up?

    fmt.Printf("Everything seems ok with %v\r\n", u)
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
    } else if u.Scheme == "" {
        remote = "gemini://" + u.Host + "/" + u.Path
    }

    return remote, nil
}
