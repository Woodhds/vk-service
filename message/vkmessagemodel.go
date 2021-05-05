package message

type VkMessageModel struct {
	ID           int        `json:"id"`
	FromID       int        `json:"fromId"`
	Date         *Timestamp `json:"date"`
	Images       []string   `json:"images"`
	LikesCount   int        `json:"likesCount"`
	Owner        string     `json:"owner"`
	OwnerID      int        `json:"ownerId"`
	RepostedFrom int        `json:"repostedFrom"`
	RepostsCount int        `json:"repostsCount"`
	Text         string     `json:"text"`
}
