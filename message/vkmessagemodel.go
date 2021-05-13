package message

import (
	"fmt"
	"strings"
)

type ImageArray []string

type VkMessageModel struct {
	ID           int        `json:"id"`
	FromID       int        `json:"fromId"`
	Date         *Timestamp `json:"date"`
	Images       ImageArray `json:"images"`
	LikesCount   int        `json:"likesCount"`
	Owner        string     `json:"owner"`
	OwnerID      int        `json:"ownerId"`
	RepostedFrom int        `json:"repostedFrom"`
	RepostsCount int        `json:"repostsCount"`
	Text         string     `json:"text"`
	UserReposted bool       `json:"userReposted"`
}

func (n *ImageArray) Scan(value interface{}) error {
	if value == nil {
		*n = make(ImageArray, 0)
	}
	s := fmt.Sprint(value)

	*n = ImageArray(strings.Split(s, ";"))

	return nil
}
