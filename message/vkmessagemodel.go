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
	RepostsCount int        `json:"repostsCount"`
	Text         string     `json:"text"`
	UserReposted bool       `json:"userReposted"`
}

func (n *ImageArray) Scan(value interface{}) error {
	if value == nil {
		*n = make(ImageArray, 0)
	}
	s := fmt.Sprint(value)

	*n = strings.Split(s, ";")

	return nil
}

func New(post *VkMessage, groups []*VkGroup) *VkMessageModel {

	likes := 0
	if post.Likes != nil {
		likes = post.Likes.Count
	}
	reposts := 0
	userReposted := false
	if post.Reposts != nil {
		reposts = post.Reposts.Count
		userReposted = post.Reposts.UserReposted == 1
	}

	model := &VkMessageModel{
		ID:           post.ID,
		FromID:       post.FromID,
		Date:         post.Date,
		Images:       []string{},
		LikesCount:   likes,
		Owner:        "",
		OwnerID:      post.OwnerID,
		RepostsCount: reposts,
		Text:         post.Text,
		UserReposted: userReposted,
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

	return model
}
