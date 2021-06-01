package vkclient

import (
	"crypto/tls"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/url"
	"time"
)

type BaseClient struct {
	client  *http.Client
	baseUrl *url.URL
}

type Response struct {
	reader io.ReadCloser
}

func (resp *Response) Read(dest *interface{}) error {
	defer resp.reader.Close()

	b, err := io.ReadAll(resp.reader)

	if err != nil {
		return err
	}

	err = json.Unmarshal(b, &dest)

	if err != nil {
		return err
	}

	return nil
}

func New(accessToken string, v string) (*BaseClient, error) {
	if accessToken == "" {
		return nil, errors.New("access token required")
	}

	if v == "" {
		return nil, errors.New("version required")
	}

	httpClient := &http.Client{
		Timeout: time.Second * 30,
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
	}

	addr, e := url.Parse("https://api.vk.com/method/")

	if e != nil {
		return nil, e
	}

	return &BaseClient{
		client:  httpClient,
		baseUrl: addr,
	}, nil
}

func (client *BaseClient) Request(method string, url string) (*Response, error) {
	request, err := http.NewRequest(method, client.baseUrl.String(), nil)

	if err != nil {
		return nil, err
	}
	response, err := client.client.Do(request)

	if err != nil {
		return nil, err
	}

	return &Response{response.Body}, nil

}
