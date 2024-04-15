package usecase_test

import (
	"errors"
	"single-window/internal/entity"
	"single-window/internal/usecase"
	"testing"

	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func initMessagesMocks(t *testing.T) (*usecase.MessageUseCase, *MockIMessageRepo) {
	t.Helper()

	mockCtl := gomock.NewController(t)
	defer mockCtl.Finish()

	messageMockRepo := NewMockIMessageRepo(mockCtl)
	messageMockUseCase := usecase.NewMessageUseCase(messageMockRepo, nil, nil)

	return messageMockUseCase, messageMockRepo
}

var msgClaims = entity.Claims{
	UserId: "11111111-AAAA-AAAA-AAAA-AAAAAAAAAAAA",
}

func TestGetMessagesV1(t *testing.T) {
	t.Parallel()

	messageUseCase, messageMockRepo := initMessagesMocks(t)

	t.Run("Empty result", func(t *testing.T) {
		messageMockRepo.EXPECT().GetMessagesV1(msgClaims.UserId, disputeId, 11, 0).Return(nil, nil)

		result, _, err := messageUseCase.GetMessagesV1(&msgClaims, disputeId, 10, 0)

		require.NoError(t, err)
		require.Nil(t, result)
	})

	t.Run("Get messages", func(t *testing.T) {
		expectedMessages := []entity.Message{{}, {}, {}}
		messageMockRepo.EXPECT().GetMessagesV1(msgClaims.UserId, disputeId, 4, 0).Return(expectedMessages, nil)

		result, _, err := messageUseCase.GetMessagesV1(&msgClaims, disputeId, 3, 0)

		require.NoError(t, err)
		require.Equal(t, expectedMessages, result)
	})

	t.Run("Failed to get messages", func(t *testing.T) {
		mockError := errors.New("mock error")
		messageMockRepo.EXPECT().GetMessagesV1(msgClaims.UserId, disputeId, 2, 0).Return(nil, mockError)

		_, _, err := messageUseCase.GetMessagesV1(&msgClaims, disputeId, 1, 0)

		require.Error(t, err)
		require.Equal(t, mockError, errors.Unwrap(err))
	})
}
