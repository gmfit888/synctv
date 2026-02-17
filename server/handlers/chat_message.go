package handlers
 
import (
	"net/http"
	"time"
 
	"github.com/gin-gonic/gin"
	"github.com/synctv-org/synctv/internal/db"
	"github.com/synctv-org/synctv/internal/model"
	"github.com/synctv-org/synctv/internal/op"
	"github.com/synctv-org/synctv/server/model"
)
 
// GetRoomChatMessages 获取房间聊天记录
func GetRoomChatMessages(ctx *gin.Context) {
	roomID := ctx.Param("room_id")
	if roomID == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "room_id is required"})
		return
	}
 
	// 验证房间是否存在且有权限访问
	_, err := db.GetRoomByID(roomID)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "room not found"})
		return
	}
 
	var req model.ChatMessagesReq
	if err := req.Decode(ctx); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
 
	// 设置默认分页参数
	if req.Page <= 0 {
		req.Page = 1
	}
	if req.PageSize <= 0 {
		req.PageSize = 50
	}
 
	messages, total, err := db.GetRoomChatMessages(roomID, req.Page, req.PageSize)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get chat messages"})
		return
	}
 
	// 转换为响应格式
	resp := &model.ChatMessagesResp{
		Messages: make([]model.ChatMessageResp, 0, len(messages)),
		Total:    total,
		Page:     req.Page,
		PageSize: req.PageSize,
	}
 
	for _, msg := range messages {
		resp.Messages = append(resp.Messages, model.ChatMessageResp{
			ID:        msg.ID,
			CreatedAt: msg.CreatedAt.UnixMilli(),
			RoomID:    msg.RoomID,
			UserID:    msg.UserID,
			Username:  msg.Username,
			Content:   msg.Content,
		})
	}
 
	ctx.JSON(http.StatusOK, resp)
}
 
// GetRoomRecentChatMessages 获取房间最近的聊天记录
func GetRoomRecentChatMessages(ctx *gin.Context) {
	roomID := ctx.Param("room_id")
	if roomID == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "room_id is required"})
		return
	}
 
	// 验证房间是否存在
	_, err := db.GetRoomByID(roomID)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "room not found"})
		return
	}
 
	// 默认获取最近100条
	messages, err := db.GetRoomRecentChatMessages(roomID, 100)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get recent chat messages"})
		return
	}
 
	// 转换为响应格式
	resp := make([]model.ChatMessageResp, 0, len(messages))
	for _, msg := range messages {
		resp = append(resp, model.ChatMessageResp{
			ID:        msg.ID,
			CreatedAt: msg.CreatedAt.UnixMilli(),
			RoomID:    msg.RoomID,
			UserID:    msg.UserID,
			Username:  msg.Username,
			Content:   msg.Content,
		})
	}
 
	ctx.JSON(http.StatusOK, resp)
}
 
// DeleteRoomChatMessages 删除房间所有聊天记录
func DeleteRoomChatMessages(ctx *gin.Context) {
	roomID := ctx.Param("room_id")
	if roomID == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "room_id is required"})
		return
	}
 
	// 验证权限
	userID := ctx.GetString("user_id")
	if userID == "" {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
 
	// 检查是否为房主或管理员
	room, err := db.GetRoomByID(roomID)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "room not found"})
		return
	}
 
	if room.CreatorID != userID {
		ctx.JSON(http.StatusForbidden, gin.H{"error": "only room creator can delete chat messages"})
		return
	}
 
	err = db.DeleteRoomChatMessages(roomID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to delete chat messages"})
		return
	}
 
	ctx.JSON(http.StatusOK, gin.H{"message": "chat messages deleted successfully"})
}
