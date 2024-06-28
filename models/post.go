package models

type Post struct {
	ID           int         `json:"id"`
	Title        string      `json:"title"`
	Content      string      `json:"content"`
	CreatedAt    string      `json:"createdAt"`
	UserID       int         `json:"userId"`
	User         UserProfile `json:"user"`
	PhotoURLs    []string    `json:"photoUrls"`
	LikedBy      []int       `json:"likedBy"`
	CommentCount int         `json:"commentCount"`
}

type CreatePostRequest struct {
	Title     string   `json:"title"`
	Content   string   `json:"content"`
	CreatedAt string   `json:"createdAt"`
	UserID    int      `json:"userId"`
	PhotoURLs []string `json:"photoUrls"`
}
