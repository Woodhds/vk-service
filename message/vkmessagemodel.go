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

func New(post *VkMessage, id int, groups []*VkGroup) *VkMessageModel {
	model := &VkMessageModel{
		ID:           post.ID,
		FromID:       post.FromID,
		Date:         post.Date,
		Images:       []string{},
		LikesCount:   post.Likes.Count,
		Owner:        "",
		OwnerID:      post.OwnerID,
		RepostedFrom: id,
		RepostsCount: post.Reposts.Count,
		Text:         post.Text,
	}

	for _, i := range post.Attachments {
		if len(i.Photo.Sizes) > 2 {
			model.Images = append(model.Images, i.Photo.Sizes[3].Url)
		}
	}

	for _, g := range groups {
		if g.ID == -post.OwnerID {
			model.Owner = g.Name
		}
	}

	if post.Reposts.UserReposted == 1 {
		model.UserReposted = true
	}

	return model
}
