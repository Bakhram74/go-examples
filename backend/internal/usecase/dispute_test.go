package usecase_test

import (
	"errors"
	"fmt"
	"single-window/internal/entity"
	"single-window/internal/usecase"
	"single-window/pkg/accessor"
	accessortest "single-window/pkg/accessor/accessor_test"
	"testing"

	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

const (
	disputeId = "AAAAAAAA-AAAA-AAAA-AAAA-AAAAAAAAAAAA"
	userId    = "ZZZZZZZZ-ZZZZ-ZZZZ-ZZZZ-ZZZZZZZZZZZZ"
)

func initMocks(t *testing.T) (*usecase.DisputeUseCase, *MockIDisputeRepo, *accessortest.MockIAccessor) {
	t.Helper()

	mockCtl := gomock.NewController(t)
	defer mockCtl.Finish()

	disputeMockRepo := NewMockIDisputeRepo(mockCtl)
	accessorMock := accessortest.NewMockIAccessor(mockCtl)
	disputeMockUseCase := usecase.NewDisputeUseCase(disputeMockRepo, nil, accessorMock)

	return disputeMockUseCase, disputeMockRepo, accessorMock
}

// TODO: Добавить проверку accessor
func TestGetDisputesV1(t *testing.T) {
	t.Parallel()

	disputeUsecase, disputeMockRepo, accessorMock := initMocks(t)

	t.Run("Empty result", func(t *testing.T) {
		var employeeId int64 = 0
		claims := entity.Claims{
			EmployeeId: &employeeId,
			UserRole:   "admin",
		}

		accessorMock.EXPECT().CheckAccess(gomock.Any(), gomock.Any(), gomock.Any()).Return(accessor.AccessGranted, nil)
		disputeMockRepo.EXPECT().GetDisputesV1(claims, 11, 0, entity.Any).Return(nil, nil)
		result, _, err := disputeUsecase.GetDisputesV1(claims, 10, 0, entity.Any)

		require.NoError(t, err)
		require.Nil(t, result)
	})

	t.Run("Get disputes", func(t *testing.T) {
		var employeeId int64 = 0
		expectedDisputes := []entity.DisputeList{{}, {}, {}}
		claims := entity.Claims{
			EmployeeId: &employeeId,
			UserRole:   "admin",
		}

		accessorMock.EXPECT().CheckAccess(gomock.Any(), gomock.Any(), gomock.Any()).Return(accessor.AccessGranted, nil)
		disputeMockRepo.EXPECT().GetDisputesV1(gomock.Any(), 4, 0, entity.Any).Return(expectedDisputes, nil)
		result, _, err := disputeUsecase.GetDisputesV1(claims, 3, 0, entity.Any)

		require.NoError(t, err)
		require.Equal(t, expectedDisputes, result)
	})

	t.Run("Error handling", func(t *testing.T) {
		var employeeId int64 = 0
		claims := entity.Claims{
			EmployeeId: &employeeId,
			UserRole:   "admin",
		}
		mockError := errors.New("mock error")

		accessorMock.EXPECT().CheckAccess(gomock.Any(), gomock.Any(), gomock.Any()).Return(accessor.AccessGranted, nil)
		disputeMockRepo.EXPECT().GetDisputesV1(gomock.Any(), 2, 0, entity.Any).Return(nil, mockError)
		_, _, err := disputeUsecase.GetDisputesV1(claims, 1, 0, entity.Any)

		require.Error(t, err)
		require.Equal(t, mockError, errors.Unwrap(err))
	})
}

func TestCreateDisputeV1(t *testing.T) {
	t.Parallel()

	disputeUsecase, disputeMockRepo, _ := initMocks(t)

	userId := "user_id"
	messageBody := "message_body"
	shortageId := "shortage_id"
	organizationId := "organization_id"

	t.Run("Dispute created", func(t *testing.T) {
		dispute := entity.CreateDispute{}
		createDisputeReqBody := entity.CreateDisputeV1MultipartBody{
			MessageBody:    messageBody,
			ShortageId:     shortageId,
			OrganizationId: organizationId,
		}

		disputeMockRepo.EXPECT().CreateDisputeV1(&createDisputeReqBody, userId, nil).Return(&dispute, nil)

		result, err := disputeUsecase.CreateDisputeV1(&entity.Claims{
			UserId: userId,
		}, &createDisputeReqBody)

		require.NoError(t, err)
		require.Equal(t, &dispute, result)
	})

	t.Run("Dispute create error (dispute already exists)", func(t *testing.T) {
		createDisputeReqBody := entity.CreateDisputeV1MultipartBody{
			MessageBody:    messageBody,
			ShortageId:     shortageId,
			OrganizationId: organizationId,
		}
		mockErr := fmt.Errorf("dispute on shartage_id: %s has already been created", createDisputeReqBody.ShortageId)

		disputeMockRepo.EXPECT().CreateDisputeV1(&createDisputeReqBody, userId, nil).Return(nil, mockErr)

		result, err := disputeUsecase.CreateDisputeV1(&entity.Claims{UserId: userId}, &createDisputeReqBody)

		require.Nil(t, result)
		require.ErrorContains(t, err, "has already been created")
	})

	t.Run("Shortage is not available for user", func(t *testing.T) {
		createDisputeReqBody := entity.CreateDisputeV1MultipartBody{
			MessageBody:    messageBody,
			ShortageId:     "not user shortage",
			OrganizationId: organizationId,
		}
		mockErr := fmt.Errorf("failed to find shortage with id = %s on user_id = %s", createDisputeReqBody.ShortageId, userId)

		disputeMockRepo.EXPECT().CreateDisputeV1(&createDisputeReqBody, userId, nil).Return(nil, mockErr)

		result, err := disputeUsecase.CreateDisputeV1(&entity.Claims{UserId: userId}, &createDisputeReqBody)

		require.Nil(t, result)
		require.ErrorContains(t, err, mockErr.Error())
	})
}

func TestGetDisputeV1(t *testing.T) {
	t.Parallel()

	disputeUsecase, disputeMockRepo, _ := initMocks(t)
	disputeId := "AAAAAAAA-AAAA-AAAA-AAAA-AAAAAAAAAAAA"

	t.Run("Has result with dispute_id and goods_id", func(t *testing.T) {
		goodsId := 9999999999

		dispute := entity.Dispute{
			DisputeId: disputeId,
			GoodsId:   1234567890,
		}

		workersList := []string{"Filin", "Sneak"}

		disputeMockRepo.EXPECT().GetGuiltyWorkersV1(disputeId).Return(workersList, nil)
		disputeMockRepo.EXPECT().GetDisputeV1(&disputeId, &goodsId).Return(&dispute, nil)

		result, err := disputeUsecase.GetDisputeV1(&disputeId, &goodsId)

		require.NoError(t, err)
		require.NotEqual(t, &dispute.GoodsId, goodsId)
		require.Equal(t, &dispute, result)
		require.Equal(t, result.GuiltyWorkerNames, workersList)
	})

	t.Run("Has result with dispute_id only", func(t *testing.T) {
		dispute := entity.Dispute{
			DisputeId: disputeId,
		}

		workersList := []string{"Filin", "Sneak"}

		disputeMockRepo.EXPECT().GetGuiltyWorkersV1(disputeId).Return(workersList, nil)
		disputeMockRepo.EXPECT().GetDisputeV1(&disputeId, nil).Return(&dispute, nil)

		result, err := disputeUsecase.GetDisputeV1(&disputeId, nil)

		require.NoError(t, err)
		require.Equal(t, &dispute, result)
		require.Equal(t, result.GuiltyWorkerNames, workersList)
	})

	t.Run("Has result with goods_id only", func(t *testing.T) {
		goodsId := 1234567890

		dispute := entity.Dispute{
			DisputeId: disputeId,
			GoodsId:   goodsId,
		}

		workersList := []string{"Filin", "Sneak"}

		disputeMockRepo.EXPECT().GetGuiltyWorkersV1(disputeId).Return(workersList, nil)
		disputeMockRepo.EXPECT().GetDisputeV1(nil, &goodsId).Return(&dispute, nil)

		result, err := disputeUsecase.GetDisputeV1(nil, &goodsId)

		require.NoError(t, err)
		require.Equal(t, &dispute, result)
		require.Equal(t, result.GuiltyWorkerNames, workersList)
	})

	t.Run("Guilty workers not found DB error", func(t *testing.T) {
		mockError := errors.New("mockError")

		dispute := entity.Dispute{
			DisputeId: disputeId,
		}

		disputeMockRepo.EXPECT().GetGuiltyWorkersV1(disputeId).Return(nil, mockError)
		disputeMockRepo.EXPECT().GetDisputeV1(&disputeId, nil).Return(&dispute, nil)

		result, err := disputeUsecase.GetDisputeV1(&disputeId, nil)

		require.ErrorIs(t, err, mockError)
		require.Equal(t, &dispute, result)
	})

	t.Run("Dispute fetch from DB error", func(t *testing.T) {
		mockError := errors.New("DB error")
		workersList := []string{"Filin", "Sneak"}

		disputeMockRepo.EXPECT().GetGuiltyWorkersV1(disputeId).Return(workersList, nil)
		disputeMockRepo.EXPECT().GetDisputeV1(&disputeId, nil).Return(nil, mockError)

		result, err := disputeUsecase.GetDisputeV1(&disputeId, nil)

		require.ErrorIs(t, err, mockError)
		require.Nil(t, result)
	})

	t.Run("Dispute error query empty", func(t *testing.T) {
		result, err := disputeUsecase.GetDisputeV1(nil, nil)

		require.Error(t, err)
		require.Nil(t, result)
	})
}

func TestCloseDisputeV1V1(t *testing.T) {
	t.Parallel()

	disputeUsecase, disputeMockRepo, accessorMock := initMocks(t)

	t.Run("Dispute close error", func(t *testing.T) {
		mockError := errors.New("close error")
		userIds := []string{"USERID-1"}

		accessorMock.EXPECT().CheckAccess(gomock.Any(), gomock.Any(), gomock.Any()).Return(accessor.AccessGranted, nil)

		disputeMockRepo.EXPECT().CloseDisputeV1(disputeId, userIds).Return(nil, mockError)

		result, err := disputeUsecase.CloseDisputeV1(&entity.Claims{}, disputeId, userIds)

		require.ErrorContains(t, err, "CloseDisputeV1 - UseCase - CloseDispute")
		require.ErrorIs(t, err, mockError)
		require.Nil(t, result)
	})

	t.Run("No arguments error", func(t *testing.T) {
		accessorMock.EXPECT().CheckAccess(gomock.Any(), gomock.Any(), gomock.Any()).Return(accessor.AccessGranted, nil)

		result, err := disputeUsecase.CloseDisputeV1(&entity.Claims{}, "", nil)

		require.ErrorContains(t, err, "error disputeId has null or guiltyWorkersIds empty")
		require.Nil(t, result)
	})

	t.Run("Guilty workers successfully inserted", func(t *testing.T) {
		var disputeClosed entity.CloseDispute
		userIds := []string{"USERID-1", "USERID-2"}
		userNames := []string{"USERNAME-1", "USERNAME-2"}

		accessorMock.EXPECT().CheckAccess(gomock.Any(), gomock.Any(), gomock.Any()).Return(accessor.AccessGranted, nil)

		disputeMockRepo.EXPECT().CloseDisputeV1(disputeId, userIds).Return(&disputeClosed, nil)
		disputeMockRepo.EXPECT().GetGuiltyWorkersV1(disputeId).Return(userNames, nil)

		result, err := disputeUsecase.CloseDisputeV1(&entity.Claims{}, disputeId, userIds)

		require.Equal(t, result, &disputeClosed)
		require.NoError(t, err)
	})

	t.Run("Error get worker names", func(t *testing.T) {
		mockError := errors.New("worker names not found")
		var disputeClosed entity.CloseDispute
		userIds := []string{"USERID-1", "USERID-2"}
		accessorMock.EXPECT().CheckAccess(gomock.Any(), gomock.Any(), gomock.Any()).Return(accessor.AccessGranted, nil)

		disputeMockRepo.EXPECT().CloseDisputeV1(disputeId, userIds).Return(&disputeClosed, nil)
		disputeMockRepo.EXPECT().GetGuiltyWorkersV1(disputeId).Return(nil, mockError)

		result, err := disputeUsecase.CloseDisputeV1(&entity.Claims{}, disputeId, userIds)

		require.ErrorIs(t, err, mockError)
		require.ErrorContains(t, err, "CloseDisputeV1 - UseCase - getGuiltyWorkers")
		require.Equal(t, result, &disputeClosed)
	})
}

func TestGetShortageV1(t *testing.T) {
	t.Parallel()

	disputeUsecase, disputeMockRepo, _ := initMocks(t)

	t.Run("Empty result", func(t *testing.T) {
		disputeMockRepo.EXPECT().GetShortagesV1(userId, 11, 0).Return(nil, nil)

		result, _, err := disputeUsecase.GetShortagesV1(userId, 10, 0)

		require.NoError(t, err)
		require.Nil(t, result)
	})

	t.Run("Get shoratges", func(t *testing.T) {
		expectedShortages := []entity.Shortage{{}, {}, {}}
		disputeMockRepo.EXPECT().GetShortagesV1(userId, 4, 0).Return(expectedShortages, nil)

		result, _, err := disputeUsecase.GetShortagesV1(userId, 3, 0)

		require.NoError(t, err)
		require.Equal(t, expectedShortages, result)
	})

	t.Run("Error handling", func(t *testing.T) {
		mockError := errors.New("mock error")
		disputeMockRepo.EXPECT().GetShortagesV1(userId, 2, 0).Return(nil, mockError)

		_, _, err := disputeUsecase.GetShortagesV1(userId, 1, 0)

		require.Error(t, err)
		require.Equal(t, mockError, errors.Unwrap(err))
	})
}
