package model
 
import (
	"time"
 
	"github.com/synctv-org/synctv/utils"
	"gorm.io/gorm"
)
 
type ChatMessage struct {
	ID        string    `gorm:"primaryKey;type:char(32)" json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	RoomID    string    `gorm:"index;type:char(32);not null" json:"room_id"`
	UserID    string    `gorm:"index;type:char(32);not null" json:"user_id"`
	Username  string    `gorm:"type:varchar(32);not null" json:"username"`
	Content   string    `gorm:"type:text;not null" json:"content"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}
 
func (c *ChatMessage) BeforeCreate(_ *gorm.DB) error {
	if c.ID == "" {
		c.ID = utils.SortUUID()
	}
	return nil
}
