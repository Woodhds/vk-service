package predictor

import (
	"context"
	vkPostPredict "github.com/woodhds/vk.service/gen/predict"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"
	"time"
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
}

type SavePredictRequest struct {
	Category string `json:"category"`
}

type predictorClient struct {
	client vkPostPredict.MessagePredictServiceClient
}

func (c *predictorClient) SaveMessage(owner int, id int, text string, ownerName string, category string) error {
	reqData := &vkPostPredict.MessageSaveRequest{
		OwnerId:   int32(owner),
		Id:        int32(id),
		Text:      text,
		Category:  category,
		OwnerName: ownerName,
	}

	if _, e := c.client.Save(context.Background(), reqData); e != nil {
		return e
	}

	return nil
}

func (c *predictorClient) Predict(messages []*PredictMessage) ([]*PredictMessage, error) {
	request := createRequest(messages)

	respData, e := c.client.Predict(context.Background(), request)

	if e != nil {
		log.Print(e)
		return messages, e
	}

	for i := 0; i < len(messages); i++ {
		h := messages[i]
		for j := 0; j < len(respData.Messages); j++ {
			r := respData.Messages[j]
			if int(r.Id) == h.Id && int(r.OwnerId) == h.OwnerId {
				r.Category = h.Category
				r.IsAccept = h.IsAccept
				r.Scores = h.Scores
				break
			}
		}
	}

	return messages, nil
}

func NewClient(host string) (Predictor, error) {
	ctx, _ := context.WithTimeout(context.Background(), 30*time.Second)
	dialContext, _ := grpc.DialContext(
		ctx,
		host,
		grpc.WithTransportCredentials(insecure.NewCredentials()))
	client := vkPostPredict.NewMessagePredictServiceClient(dialContext)
	return &predictorClient{
		client: client,
	}, nil
}

func createRequest(messages []*PredictMessage) *vkPostPredict.MessagePredictRequest {
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
