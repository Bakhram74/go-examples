package usecase

import (
	"fmt"
	"single-window/internal/entity"
	"single-window/pkg/accessor"
	"single-window/pkg/s3"
	"single-window/pkg/utils"
)

type MessageUseCase struct {
	repo     IMessageRepo
	s3Client s3.IS3Client
	acc      accessor.IAccessor
}

func NewMessageUseCase(mr IMessageRepo, s3Client s3.IS3Client, acc accessor.IAccessor) *MessageUseCase {
	return &MessageUseCase{
		repo:     mr,
		s3Client: s3Client,
		acc:      acc,
	}
}

func (muc *MessageUseCase) GetMessagesV1(claims *entity.Claims, disputeId string, limit, offset int) ([]entity.Message, *bool, error) {
	messages, err := muc.repo.GetMessagesV1(claims.UserId, disputeId, limit+1, offset)
	if err != nil {
		return nil, nil, fmt.Errorf("GetMessagesV1 - UseCase %w", err)
	}

	messagesLen := len(messages)
	hasMoreData := messagesLen > limit
	if hasMoreData {
		messages = messages[:messagesLen-1]
	}

	return messages, &hasMoreData, nil
}

func (muc *MessageUseCase) CreateMessageV1(claims *entity.Claims, createMsgReqBody entity.CreateMessageV1MultipartBody) (*entity.CreateMsgResponse, error) {
	attachmentPath, err := utils.UploadFile(muc.s3Client, claims.UserId, createMsgReqBody.File)
	if err != nil {
		return nil, fmt.Errorf("failed to upload file to s3, err: %w", err)
	}

	createdMsg, err := muc.repo.CreateMessageV1(createMsgReqBody, claims.UserId, attachmentPath)
	if err != nil {
		// TODO: если ошибка, надо удалить загруженное вложение или пометить его на удаление через какое-то время
		return nil, fmt.Errorf("CreateMessageV1 - UseCase %w", err)
	}

	return createdMsg, nil
}
