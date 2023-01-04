package vkclient

import (
	"bytes"
	"fmt"
	"net/url"
	"strings"

	"github.com/woodhds/vk.service/message"
)

type WallClient interface {
	Get(request *WallGetRequest) (*message.VkWallResponse, error)
	GetById(messages []*message.VkRepostMessage, fields ...string) (*message.VkResponse, error)
}

type wallClient struct {
	baseClient *BaseClient
}

type WallGetRequest struct {
	Filter   string
	OwnerId  int
	Offset   int
	Count    int
	Extended bool
}

type RepostResponse struct {
	Success int `json:"success"`
}

func NewWallClient(token string, version string) (*wallClient, error) {
	baseClient, e := New(token, version)

	if e != nil {
		return nil, e
	}

	return &wallClient{
		baseClient: baseClient,
	}, nil
}

func (wallClient *wallClient) Get(request *WallGetRequest) (*message.VkWallResponse, error) {
	u := url.URL{}
	query := u.Query()
	query.Add("filter", request.Filter)
	query.Add("owner_id", fmt.Sprintf("%d", request.OwnerId))
	query.Add("offset", fmt.Sprintf("%d", request.Offset))
	count := request.Count
	if count <= 0 {
		count = 20
	}
	query.Add("count", fmt.Sprintf("%d", count))
	extended := 0
	if request.Extended {
		extended = 1
	}

	query.Add("extended", fmt.Sprintf("%d", extended))
	u.RawQuery = query.Encode()

	resp, e := wallClient.baseClient.Get("wall.get", u.String())

	if e != nil {
		return nil, e
	}

	data := &message.VkWallResponse{}

	e = resp.Read(data)

	if e != nil {
		return nil, e
	}

	return data, nil
}

func (wallClient *wallClient) GetById(messages []*message.VkRepostMessage, fields ...string) (*message.VkResponse, error) {
	var yt bytes.Buffer
	for _, dataItem := range messages {

		yt.WriteString(fmt.Sprintf("%d_%d,", dataItem.OwnerID, dataItem.ID))
	}

	u := fmt.Sprintf(`posts=%s&extended=1&fields=%s`, yt.String(), strings.Join(fields, ","))

	response, e := wallClient.baseClient.Get("wall.getById", u)

	if e != nil {
		return nil, e
	}

	data := &message.VkResponse{}

	if e := response.Read(data); e != nil {
		return nil, e
	}

	return data, nil
}

func (wallClient *wallClient) Repost(message *message.VkRepostMessage) error {
	q := fmt.Sprintf("object=wall%d_%d", message.OwnerID, message.ID)

	resp, e := wallClient.baseClient.Get("wall.repost", q)

	if e != nil {
		return e
	}

	var data RepostResponse

	if e := resp.Read(&data); e != nil {
		return e
	}

	return nil
}
