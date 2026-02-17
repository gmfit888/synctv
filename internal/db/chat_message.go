package db
 
import (
	"github.com/synctv-org/synctv/internal/model"
	"gorm.io/gorm"
)
 
const (
	ErrChatMessageNotFound = "chat message"
)
 
// CreateChatMessage 创建聊天记录
func CreateChatMessage(roomID, userID, username, content string) (*model.ChatMessage, error) {
	chatMessage := &model.ChatMessage{
		RoomID:   roomID,
		UserID:   userID,
		Username: username,
		Content:  content,
	}
 
	err := DB().Create(chatMessage).Error
	if err != nil {
		return nil, err
	}
 
	return chatMessage, nil
}
 
// GetRoomChatMessages 获取房间的聊天记录
func GetRoomChatMessages(roomID string, page, pageSize int) ([]*model.ChatMessage, int64, error) {
	var messages []*model.ChatMessage
	var total int64
 
	query := DB().Model(&model.ChatMessage{}).Where("room_id = ?", roomID)
	
	// 获取总数
	err := query.Count(&total).Error
	if err != nil {
		return nil, 0, err
	}
 
	// 分页查询
	offset := (page - 1) * pageSize
	err = query.Order("created_at DESC").
		Offset(offset).
		Limit(pageSize).
		Find(&messages).Error
	
	if err != nil {
		return nil, 0, err
	}
 
	// 反转消息顺序，使旧消息在前
	for i, j := 0, len(messages)-1; i < j; i, j = i+1, j-1 {
		messages[i], messages[j] = messages[j], messages[i]
	}
 
	return messages, total, nil
}
 
// GetRoomRecentChatMessages 获取房间最近的聊天记录（默认最近100条）
func GetRoomRecentChatMessages(roomID string, limit int) ([]*model.ChatMessage, error) {
	if limit <= 0 {
		limit = 100
	}
 
	var messages []*model.ChatMessage
	err := DB().Where("room_id = ?", roomID).
		Order("created_at DESC").
		Limit(limit).
		Find(&messages).Error
	
	if err != nil {
		return nil, err
	}
 
	// 反转消息顺序，使旧消息在前
	for i, j := 0, len(messages)-1; i < j; i, j = i+1, j-1 {
		messages[i], messages[j] = messages[j], messages[i]
	}
 
	return messages, nil
}
 
// DeleteRoomChatMessages 删除房间所有聊天记录（软删除）
func DeleteRoomChatMessages(roomID string) error {
	return DB().Where("room_id = ?", roomID).Delete(&model.ChatMessage{}).Error
}
 
// DeleteUserChatMessages 删除用户在房间的所有聊天记录（软删除）
func DeleteUserChatMessages(roomID, userID string) error {
	return DB().Where("room_id = ? AND user_id = ?", roomID, userID).Delete(&model.ChatMessage{}).Error
}
