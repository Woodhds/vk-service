package vkclient

import (
	"bytes"
	"fmt"
	"net/url"
	"strings"

	"github.com/woodhds/vk.service/message"
)

type WallClient struct {
	baseclient *BaseClient
}

const FILTER_OWNER = "owner"
const FILTER_ALL = "all"

type WallGetRequest struct {
	Filter   string
	OwnerId  int
	Offset   int
	Count    int
	Extended bool
}

type RepostResponse struct {
	Response struct {
		Success int `json:"success"`
	} `json:"response"`
}

func NewWallClient(token string, version string) (*WallClient, error) {
	baseClient, e := New(token, version)

	if e != nil {
		return nil, e
	}

	return &WallClient{
		baseclient: baseClient,
	}, nil
}

func (wallClient *WallClient) Get(request *WallGetRequest) (*message.VkWallResponse, error) {
	url := url.URL{}
	query := url.Query()
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
	url.RawQuery = query.Encode()

	resp, e := wallClient.baseclient.Get("wall.get", url.String())

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

func (wallClient *WallClient) GetById(messages []*message.VkRepostMessage, fields ...string) (*message.VkResponse, error) {
	var yt bytes.Buffer
	for _, dataItem := range messages {

		yt.WriteString(fmt.Sprintf("%d_%d,", dataItem.OwnerID, dataItem.ID))
	}

	url := fmt.Sprintf(`posts=%s&extended=1&fields=%s`, yt.String(), strings.Join(fields, ","))

	response, e := wallClient.baseclient.Get("wall.getById", url)

	if e != nil {
		return nil, e
	}

	data := &message.VkResponse{}

	response.Read(data)

	return data, nil
}

func (walllClient *WallClient) Repost(message *message.VkRepostMessage) error {
	q := fmt.Sprintf("object=wall%d_%d", message.OwnerID, message.ID)

	resp, e := walllClient.baseclient.Get("wall.repost", q)

	if e != nil {
		return e
	}

	var data RepostResponse

	if e := resp.Read(&data); e != nil {
		return e
	}

	return nil
}
