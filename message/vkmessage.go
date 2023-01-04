package message

import (
	"fmt"
	"strconv"
	"time"
)

type Timestamp time.Time

type VkRepostMessage struct {
	OwnerID int `json:"owner_id"`
	ID      int `json:"id"`
}

type SimpleMessageModel struct {
	OwnerID int    `json:"ownerId"`
	ID      int    `json:"id"`
	Text    string `json:"text"`
}

type VkWallResponse struct {
	Items []struct {
		*VkMessage
		CopyHistory []struct {
			OwnerID int `json:"owner_id"`
			ID      int `json:"id"`
		} `json:"copy_history"`
	} `json:"items"`
}

type VkResponse struct {
	Items    []*VkMessage `json:"items"`
	Groups   []*VkGroup   `json:"groups"`
	Profiles []*VkProfile `json:"profiles"`
}

type VkMessage struct {
	Date        *Timestamp `json:"date"`
	OwnerID     int        `json:"owner_id"`
	ID          int        `json:"id"`
	Text        string     `json:"text"`
	FromID      int        `json:"from_id"`
	Reposts     *VkReposts `json:"reposts"`
	Likes       *VkLikes   `json:"likes"`
	Attachments []struct {
		Photo struct {
			Sizes []struct {
				Url  string `json:"url"`
				Type string `json:"type"`
			} `json:"sizes"`
		} `json:"photo"`
	} `json:"attachments"`
}

type VkGroup struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	IsMember int    `json:"is_member"`
}

type VkReposts struct {
	Count        int `json:"count"`
	UserReposted int `json:"user_reposted"`
}

type VkLikes struct {
	Count int `json:"count"`
}

type VkProfile struct {
	ID        int    `json:"id"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
}

func (t *Timestamp) Time() time.Time {
	return time.Time(*t)
}

func (t *Timestamp) MarshalJSON() ([]byte, error) {
	stamp := fmt.Sprintf("\"%s\"", time.Time(*t).UTC().Format(time.RFC3339))
	return []byte(stamp), nil
}

func (t *Timestamp) UnmarshalJSON(b []byte) error {
	ts, err := strconv.Atoi(string(b))
	if err != nil {
		return err
	}

	*t = Timestamp(time.Unix(int64(ts), 0))

	return nil
}
