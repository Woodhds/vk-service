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

	var data int

	e = resp.Read(&data)

	if e != nil || data == 0 {
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
	var data []*message.VkGroup

	e = resp.Read(&data)

	if e != nil || len(data) == 0 {
		return nil, errors.New("error join")
	}

	return data, nil
}
