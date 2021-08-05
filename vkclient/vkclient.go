package vkclient

import (
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"sync"
	"time"
)

type BaseClient struct {
	client      *http.Client
	m           *sync.Mutex
	baseUrl     *url.URL
	accessToken string
	version     string
}

type Response struct {
	reader *http.Response
}

var mutex sync.Mutex

func (resp *Response) Read(dest interface{}) error {
	defer resp.reader.Body.Close()

	b, err := io.ReadAll(resp.reader.Body)

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
		client:      httpClient,
		baseUrl:     addr,
		accessToken: accessToken,
		version:     v,
		m:           &mutex,
	}, nil
}

func (client *BaseClient) Request(httpMethod string, path string, query string, body io.Reader) (*Response, error) {
	u := *client.baseUrl
	u.Path = fmt.Sprintf("%s%s", u.Path, path)
	q, _ := url.ParseQuery(query)

	q.Add("access_token", client.accessToken)
	q.Add("v", client.version)

	u.RawQuery = q.Encode()

	request, err := http.NewRequest(httpMethod, u.String(), body)

	if err != nil {
		return nil, err
	}

	client.m.Lock()

	response, err := client.client.Do(request)

	time.Sleep(time.Millisecond * 500)
	client.m.Unlock()

	if err != nil {
		return nil, err
	}

	return &Response{response}, nil

}

func (client *BaseClient) Get(method string, query string) (*Response, error) {
	return client.Request(http.MethodGet, method, query, nil)
}
