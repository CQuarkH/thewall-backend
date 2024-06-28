package models

type Comment struct {
	ID        int    `json:"id"`
	Content   string `json:"content"`
	CreatedAt string `json:"createdAt"`
	UserID    int    `json:"userId"`
	PostID    int    `json:"postId"`
}

type CommentWithUserResponse struct {
	ID        int    `json:"id"`
	Content   string `json:"content"`
	CreatedAt string `json:"createdAt"`
	UserID    int    `json:"userId"`
	PostID    int    `json:"postId"`
	Username  string `json:"username"`
	UserPhoto string `json:"userPhoto"`
}

type CreateCommentRequest struct {
	Content   string `json:"content"`
	CreatedAt string `json:"createdAt"`
	UserID    int    `json:"userId"`
	PostID    int    `json:"postId"`
}
