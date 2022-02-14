package vkclient

import (
	"errors"
	"fmt"
)

type GroupClient struct {
	baseClient *BaseClient
}

type GroupJoinResponse struct {
	Response int `json:"response"`
}

func NewGroupClient(token *string, v *string) (*GroupClient, error) {
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
