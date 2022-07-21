package predictor

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/golang/protobuf/jsonpb"
	vkPostPredict "github.com/woodhds/vk.service/gen/predict"
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

type Predictor interface {
	Predict(messages []*PredictMessage) ([]*PredictMessage, error)
	SaveMessage(owner int, id int, text string, ownerName string, category string) error
}

type SavePredictRequest struct {
	Category string `json:"category"`
}

type predictorClient struct {
	httpClient *http.Client
	host       string
}

func (c *predictorClient) SaveMessage(owner int, id int, text string, ownerName string, category string) error {
	reqData := &vkPostPredict.MessageSaveRequest{
		OwnerId:   int32(owner),
		Id:        int32(id),
		Text:      text,
		Category:  category,
		OwnerName: ownerName,
	}
	marshaler := jsonpb.Marshaler{}
	b, _ := marshaler.MarshalToString(reqData)

	if req, e := makeRequest(http.MethodPut, c.host, "predict", []byte(b)); e == nil {
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

func (c *predictorClient) Predict(messages []*PredictMessage) ([]*PredictMessage, error) {
	marshaler := jsonpb.Marshaler{}
	b, _ := marshaler.MarshalToString(makePayload(messages))

	if req, e := makeRequest(http.MethodPost, c.host, "predict", []byte(b)); e == nil {
		if resp, e := c.httpClient.Do(req); e != nil {
			return messages, e
		} else {
			if resp.StatusCode != http.StatusOK {
				return messages, errors.New(fmt.Sprintf("Server responded with status %d", resp.StatusCode))
			}

			var respData vkPostPredict.MessagePredictResponse
			jsonpb.Unmarshal(resp.Body, &respData)

			for _, r := range messages {
				for _, h := range respData.Messages {
					if int32(r.Id) == h.Id && int32(r.OwnerId) == h.OwnerId {
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

func makePayload(messages []*PredictMessage) *vkPostPredict.MessagePredictRequest {
	request := vkPostPredict.MessagePredictRequest{Messages: make([]*vkPostPredict.MessagePredictRequest_PredictRequest, len(messages))}
	for i := 0; i < len(messages); i++ {
		request.Messages[i] = &vkPostPredict.MessagePredictRequest_PredictRequest{
			OwnerId: int32(messages[i].OwnerId),
			Id:      int32(messages[i].Id),
			Text:    messages[i].Text,
		}

	}
	return &request
}
