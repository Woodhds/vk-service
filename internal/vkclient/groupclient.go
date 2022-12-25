package vkclient

import (
	"errors"
	"fmt"
	"github.com/woodhds/vk.service/message"
	"strconv"
	"strings"
)

type GroupClient struct {
	baseClient *BaseClient
}

type GroupJoinResponse struct {
	Response int `json:"response"`
}

type getGroupResponse struct {
	Response []*message.VkGroup `json:"response"`
}

func NewGroupClient(token string, v string) (*GroupClient, error) {
	client, err := New(token, v)
	if err != nil {
		return nil, err
	}

	return &GroupClient{
		baseClient: client,
	}, nil
}

func (groupClient *GroupClient) Join(groupId int) error {
	resp, e := groupClient.baseClient.Get("groups.join", fmt.Sprintf("group_id=%d", groupId))

	if e != nil {
		return e
	}

	var data GroupJoinResponse

	e = resp.Read(&data)

	if e != nil || data.Response == 0 {
		return errors.New("error join")
	}

	return nil
}

func (groupClient *GroupClient) Get(groupIds []int) ([]*message.VkGroup, error) {
	if len(groupIds) > 200 {
		return nil, errors.New("too many ids")
	}

	s := make([]string, len(groupIds), len(groupIds))

	for i := 0; i < len(groupIds); i++ {
		s[i] = strconv.Itoa(groupIds[i])
	}

	request := strings.Join(s, ",")

	resp, e := groupClient.baseClient.Get("groups.getById", request)

	if e != nil {
		return nil, e
	}
	var data getGroupResponse

	e = resp.Read(&data)

	if e != nil || len(data.Response) == 0 {
		return nil, errors.New("error join")
	}

	return data.Response, nil
}
