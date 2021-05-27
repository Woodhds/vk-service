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

type VkClient interface {
	New(accessToken string, v string) (error, interface{})
	Request(method string, url string, body io.Reader)
}

type BaseClient struct {
	client  *http.Client
	baseUrl *url.URL
}

type Response struct {
}

func (client *BaseClient) New(accessToken string, v string) (*BaseClient, error) {
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

func (client *BaseClient) Request(method string, url string) (*interface{}, error) {
	request, err := http.NewRequest(method, client.baseUrl.String(), nil)

	if err != nil {
		return nil, err
	}
	response, err := client.client.Do(request)

	if err != nil {
		return nil, err
	}

	defer response.Body.Close()

	b, err := io.ReadAll(request.Body)

	if err != nil {
		return nil, err
	}

	data := &Response{}

	err = json.Unmarshal(b, &data)

	if err != nil {
		return nil, err
	}

	return data, nil
}
