package repo

import (
	"errors"
	"fmt"
	"single-window/internal/entity"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type MessageRepo struct {
	db *gorm.DB
}

func NewMessageRepo(db *gorm.DB) *MessageRepo {
	return &MessageRepo{db}
}

func (m *MessageRepo) GetMessagesV1(userId, disputeId string, limit, offset int) ([]entity.Message, error) {
	var messages []entity.Message
	err := m.db.Table("disputes.chat AS dc").
		Select("dc.*, s.role, s.avatar_url AS sender_avatar_url, s.fullname AS sender_name").
		Joins("JOIN disputes.dispute_role AS dr ON dc.dispute_id = dr.dispute_id").
		Joins("JOIN disputes.dispute as d ON dc.dispute_id = d.dispute_id").
		Joins("JOIN subjects.user AS s ON dc.sender_id = s.user_id").
		Where("dc.dispute_id = ? AND (dr.user_id = ? OR d.status = 'opened')", disputeId, userId).
		Limit(limit).
		Offset(offset).
		Order("dc.created_at DESC").
		Find(&messages).Error
	if err != nil {
		return nil, fmt.Errorf("failed to get chats via dispute_id, err: %w", err)
	}

	return messages, nil
}

func (m *MessageRepo) CreateMessageV1(createMsgReqBody entity.CreateMessageV1MultipartBody, userId string, attachPath *string) (*entity.CreateMsgResponse, error) {
	err := m.db.Table("disputes.dispute_role").
		Where("dispute_id = ? AND user_id = ?", createMsgReqBody.DisputeId, userId).First(&entity.DisputeRole{}).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("user is not a participant of the dispute")
		}
		return nil, fmt.Errorf("failed to check user participation in dispute: %w", err)
	}

	message := entity.DisputeChat{
		MessageId:      uuid.NewString(),
		SenderId:       userId,
		DisputeId:      createMsgReqBody.DisputeId,
		MessageBody:    createMsgReqBody.MessageBody,
		AttachmentPath: attachPath,
		CreatedAt:      time.Now(),
	}

	err = m.db.Table(disputesChatsTableName).Create(&message).Error
	if err != nil {
		return nil, fmt.Errorf("failed to create msg, err: %w", err)
	}

	return &entity.CreateMsgResponse{MessageId: message.MessageId}, nil
}
