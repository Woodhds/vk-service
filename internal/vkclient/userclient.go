package vkclient

import (
	"fmt"
)

type VkUserMdodel struct {
	Id     int    `json:"id"`
	Name   string `json:"name"`
	Avatar string `json:"avatar"`
}

type UserClient struct {
	baseClient *BaseClient
}

type vkUserResponse struct {
	Response struct {
		Items []struct {
			Id        int    `json:"id"`
			FirstName string `json:"first_name"`
			LastName  string `json:"last_name"`
			Photo50   string `json:"photo_50"`
		} `json:"items"`
	} `json:"response"`
}

func NewUserClient(token string, v string) (*UserClient, error) {
	client, err := New(token, v)
	if err != nil {
		return nil, err
	}

	return &UserClient{
		baseClient: client,
	}, nil
}

func (userClient *UserClient) Search(q string) ([]*VkUserMdodel, error) {
	resp, err := userClient.baseClient.Get("users.search", fmt.Sprintf("q=%s&fields=photo_50", q))
	if err != nil {
		return nil, err
	}

	data := &vkUserResponse{}

	resp.Read(data)

	result := make([]*VkUserMdodel, len(data.Response.Items))

	for i, u := range data.Response.Items {
		result[i] = &VkUserMdodel{u.Id, fmt.Sprintf("%s %s", u.FirstName, u.LastName), u.Photo50}
	}

	return result, nil
}
