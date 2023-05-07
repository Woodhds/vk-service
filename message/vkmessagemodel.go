package message

import (
	"context"
	"database/sql"
	pb "github.com/woodhds/vk.service/gen/messages"
	"google.golang.org/protobuf/types/known/timestamppb"
	"time"
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

type GroupModel struct {
	Id     int64  `json:"id"`
	Name   string `json:"name"`
	Avatar string `json:"avatar"`
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
		if len(i.Photo.Sizes) > 3 {
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

func (m *VkMessageModel) Save(conn *sql.Conn, ctx context.Context) error {
	_, sqlErr := conn.ExecContext(ctx,
		`
insert into messages (Id, FromId, Date, OwnerId, Text) 
values ($1, $2, $3, $4, $5) 
ON CONFLICT(id, ownerId) DO NOTHING`,
		m.ID, m.FromID, time.Time(*m.Date), m.OwnerID, m.Text)

	return sqlErr
}

func (m *VkMessageModel) ToDto() *pb.VkMessageExt {
	return &pb.VkMessageExt{
		Id:           int32(m.ID),
		FromId:       int32(m.FromID),
		Date:         timestamppb.New(time.Time(*m.Date)),
		Images:       m.Images,
		LikesCount:   int32(m.LikesCount),
		Owner:        m.Owner,
		OwnerId:      int32(m.OwnerID),
		RepostsCount: int32(m.RepostsCount),
		Text:         m.Text,
		UserReposted: m.UserReposted,
	}
}
