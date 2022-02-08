package predictor

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
)

type PredictMessage struct {
	OwnerId  int                `json:"ownerId"`
	Id       int                `json:"id"`
	Category string             `json:"category"`
	Text     string             `json:"text"`
	IsAccept bool               `json:"isAccept"`
	Scores   map[string]float32 `json:"scores"`
}

type PredictMessageResponse struct {
	Messages []*PredictMessage `json:"messages"`
}

type Predictor interface {
	Predict(messages []*PredictMessage) ([]*PredictMessage, error)
	SaveMessage(owner int, id int, text string, ownerName string, category string) error
	PredictMessage(message *PredictMessage) (map[string]float32, error)
}

type SavePredictRequest struct {
	Category string `json:"category"`
}

type saveRequest struct {
	OwnerId   int    `json:"ownerId"`
	Id        int    `json:"id"`
	Text      string `json:"text"`
	Category  string `json:"category"`
	OwnerName string `json:"ownerName"`
}

type predictorClient struct {
	httpClient *http.Client
	host       string
}

func (c predictorClient) SaveMessage(owner int, id int, text string, ownerName string, category string) error {
	reqData := &saveRequest{
		OwnerId:   owner,
		Id:        id,
		Text:      text,
		Category:  category,
		OwnerName: ownerName,
	}
	b, _ := json.Marshal(reqData)

	if req, e := makeRequest(http.MethodPut, c.host, "predict", b); e == nil {
		if resp, e := c.httpClient.Do(req); e != nil {
			return e
		} else {
			if resp.StatusCode != http.StatusOK {
				return errors.New(fmt.Sprintf("Server responded with status code %d", resp.StatusCode))
			}
		}
	} else {
		return e
	}

	return nil
}

func (c predictorClient) Predict(messages []*PredictMessage) ([]*PredictMessage, error) {
	b, _ := json.Marshal(messages)

	if req, e := makeRequest(http.MethodPost, c.host, "predict", b); e == nil {
		if resp, e := c.httpClient.Do(req); e != nil {
			return messages, e
		} else {
			if resp.StatusCode != http.StatusOK {
				return messages, errors.New(fmt.Sprintf("Server responded with status %d", resp.StatusCode))
			}

			var respData PredictMessageResponse

			json.NewDecoder(resp.Body).Decode(&respData)

			for _, r := range messages {
				for _, h := range respData.Messages {
					if r.Id == h.Id && r.OwnerId == h.OwnerId {
						r.Category = h.Category
						r.IsAccept = h.IsAccept
						r.Scores = h.Scores
						break
					}
				}
			}
		}

	} else {
		return messages, e
	}

	return messages, nil
}

func (c *predictorClient) PredictMessage(message *PredictMessage) (map[string]float32, error) {
	d := map[string]string{"text": message.Text}
	b, _ := json.Marshal(d)
	if req, e := makeRequest(http.MethodPost, c.host, fmt.Sprintf("predict/%d/%d", message.OwnerId, message.Id), b); e == nil {
		if resp, e := c.httpClient.Do(req); e != nil {
			return make(map[string]float32), e
		} else {
			if resp.StatusCode != http.StatusOK {
				return make(map[string]float32), errors.New(fmt.Sprintf("Server responded with status %d", resp.StatusCode))
			}

			var respData map[string]float32

			json.NewDecoder(resp.Body).Decode(&respData)

			return respData, nil
		}

	} else {
		return nil, e
	}
}

func NewClient(host string) (Predictor, error) {
	return &predictorClient{
		httpClient: &http.Client{},
		host:       host,
	}, nil
}

func makeRequest(method string, host string, path string, body []byte) (*http.Request, error) {
	u, _ := url.Parse(host)
	u.Path = path
	if req, e := http.NewRequest(method, u.String(), bytes.NewBuffer(body)); e == nil {
		req.Header.Add("content-type", "application/json")

		return req, nil
	} else {
		return nil, e
	}
}
