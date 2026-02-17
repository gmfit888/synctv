package model
 
import (
	"errors"
 
	"github.com/gin-gonic/gin"
)
 
type ChatMessagesReq struct {
	Page     int `form:"page" binding:"omitempty,min=1"`
	PageSize int `form:"page_size" binding:"omitempty,min=1,max=1000"`
}
 
func (c *ChatMessagesReq) Decode(ctx *gin.Context) error {
	return ctx.ShouldBindQuery(c)
}
 
type ChatMessageResp struct {
	ID        string `json:"id"`
	CreatedAt int64  `json:"created_at"`
	RoomID    string `json:"room_id"`
	UserID    string `json:"user_id"`
	Username  string `json:"username"`
	Content   string `json:"content"`
}
 
type ChatMessagesResp struct {
	Messages []ChatMessageResp `json:"messages"`
	Total    int64             `json:"total"`
	Page     int               `json:"page"`
	PageSize int               `json:"page_size"`
}
