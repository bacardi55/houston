package main

import (
	"context"
	"fmt"
	"os"
	"strings"
	"time"

	"net/url"

	"git.sr.ht/~adnano/go-gemini"
)

const invalidCertErrorMsg = "Certificate is unvalid"

// Nota: Gemini servers might have max execution time for cgi scripts.
// Eg: Gemserv has a 5s maximum policy before killing the request.
const maxRequestTime = 4

func main() {
	if os.Getenv("QUERY_STRING") == "" {
		fmt.Println("10 Enter the URL to test:\r\n")
		os.Exit(0)
	} else {
		fmt.Println("20 text/gemini\r\n")
	}

	remoteUrl := os.Getenv("QUERY_STRING")

	u, e := validateUrl(remoteUrl)
	if e != nil {
		fmt.Printf("Url is not valid.\r\n")
		endResponse()
	}

	invalidCert := false
	resp, err := fetchGeminiPage(u)

	if err != nil {
		if err.Error() == invalidCertErrorMsg {
			invalidCert = true
		} else {
			fmt.Printf("The url you are testing (%v) seems down.\r\n", u)
			endResponse()
		}
	}

	// if the server response if an error code, capsule isn't up.
	if resp.Status == gemini.StatusTemporaryFailure || resp.Status == gemini.StatusServerUnavailable || resp.Status == gemini.StatusCGIError || resp.Status == gemini.StatusProxyError || resp.Status == gemini.StatusPermanentFailure || resp.Status == gemini.StatusGone || resp.Status == gemini.StatusProxyRequestRefused || resp.Status == gemini.StatusBadRequest {
		fmt.Printf("The url you are testing (%v) seems down, status is %v\r\n", u, resp.Status)
		endResponse()
	}

	// Todo: Should I follow redirect and check redirect status' code…?
	// Or is a redirect response enough to see that the capsule is up?

	if invalidCert {
		fmt.Printf("The capsule %v seems up and running but the certificate is expired.\r\n", u)
	} else {
		fmt.Printf("Everything seems ok with %v.\r\n", u)
	}

	respCert := resp.TLS().PeerCertificates
	if respCert[0].PublicKeyAlgorithm.String() == "Ed25519" {
		fmt.Println("\nPS: This capsule TLS certificate is using Ed25519 algorythm that is not supported by every client.\n")
	}

	endResponse()
}

func fetchGeminiPage(remoteUrl string) (*gemini.Response, error) {
	gemclient := &gemini.Client{}
	ctx, _ := context.WithTimeout(context.Background(), time.Duration(maxRequestTime)*time.Second)
	response, err := gemclient.Get(ctx, remoteUrl)

	if err != nil {
		return response, err
	}

	if respCert := response.TLS().PeerCertificates; len(respCert) > 0 && time.Now().After(respCert[0].NotAfter) {
		return response, fmt.Errorf(invalidCertErrorMsg)
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
}

func endResponse() {
	fmt.Printf("=> /index.gmi Return to Houston homepage\r\n")
	fmt.Printf("=> /cgi-bin/houston Test another URL\r\n")
	os.Exit(0)
}
