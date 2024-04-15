package usecase_test

import (
	"errors"
	"single-window/internal/entity"
	"single-window/internal/usecase"
	"single-window/pkg/accessor"
	accessortest "single-window/pkg/accessor/accessor_test"
	"testing"

	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func initRevisionMocks(t *testing.T) (*usecase.RevisionUseCase, *MockIRevisionRepo, *accessortest.MockIAccessor) {
	t.Helper()

	mockCtl := gomock.NewController(t)
	defer mockCtl.Finish()

	messageMockRepo := NewMockIRevisionRepo(mockCtl)
	accessorMock := accessortest.NewMockIAccessor(mockCtl)
	messageMockUseCase := usecase.NewRevisionUseCase(messageMockRepo, accessorMock, nil)

	return messageMockUseCase, messageMockRepo, accessorMock
}

var revClaims = entity.Claims{
	UserId:         "11111111-AAAA-AAAA-AAAA-AAAAAAAAAAAA",
	OrganizationId: "11111111-AAAA-AAAA-AAAA-AAAAAAAAAAAA",
}

var defPagination = entity.Pagination{
	Limit:  3,
	Offset: 0,
}


// TODO: Дописать тесты по созданию ревизии
func TestCreateRevisionV1(t *testing.T) {
	t.Parallel()

	uc, repo, accessorMock := initRevisionMocks(t)

	t.Run("Revision created", func(t *testing.T) {
		revision := entity.Revision{}
		body := entity.CreateRevisionV1MultipartBody{
			DisputeId:      "22222222-AAAA-AAAA-AAAA-AAAAAAAAAAAA",
			OrganizationId: "33333333-AAAA-AAAA-AAAA-AAAAAAAAAAAA",
			MessageBody:    "MESSAGE",
		}

		accessorMock.EXPECT().CheckAccess(gomock.Any(), gomock.Any(), gomock.Any()).Return(accessor.AccessGranted, nil)

		repo.EXPECT().CreateRevisionV1(&body, revClaims.UserId, nil).Return(&revision, nil)

		result, err := uc.CreateRevisionV1(&revClaims, &body)

		require.NoError(t, err)
		require.Equal(t, result, &revision)
	})
}

// TODO: Дописать тесты по получению списка корреспонденций
func TestGetCorrespondencesV1(t *testing.T) {
	t.Parallel()

	uc, repo, accessorMock := initRevisionMocks(t)

	pag := entity.Pagination{Offset: 0, Limit: 10}

	revClaims := entity.GetCorrespondencesData{
		Claims:     revClaims,
		RevisionID: "111111111-AAAA-AAAA-AAAA-AAAAAAAAAAAA",
	}

	t.Run("Correspondences received", func(t *testing.T) {
		listCorrespondences := []entity.Correspondence{{}, {}, {}}
		accessorMock.EXPECT().CheckAccess(gomock.Any(), gomock.Any(), gomock.Any()).Return(accessor.AccessGranted, nil)

		repo.EXPECT().GetCorrespondencesV1(revClaims.RevisionID, pag.IncrementLimit()).Return(listCorrespondences, nil)

		result, hasMoreData, err := uc.GetCorrespondencesV1(revClaims, pag)

		require.NoError(t, err)
		require.Equal(t, *hasMoreData, false)
		require.Equal(t, result, listCorrespondences)
	})
}

func TestGetRevisionsV1(t *testing.T) {
	t.Parallel()

	uc, repo, accessorMock := initRevisionMocks(t)

	t.Run("Empty result", func(t *testing.T) {
		accessorMock.EXPECT().CheckAccess(gomock.Any(), gomock.Any(), gomock.Any()).Return(accessor.AccessGranted, nil)
		repo.EXPECT().
			GetRevisionsV1(entity.GetRevisionsData{Claims: revClaims}, defPagination.IncrementLimit()).
			Return(nil, nil)

		result, err := uc.GetRevisionsV1(entity.GetRevisionsData{Claims: revClaims}, defPagination)

		require.NoError(t, err)
		require.Nil(t, result.Data)
	})

	t.Run("Get revisions", func(t *testing.T) {
		expected := []entity.Revision{{}, {}, {}}
		accessorMock.EXPECT().CheckAccess(gomock.Any(), gomock.Any(), gomock.Any()).Return(accessor.AccessGranted, nil)
		repo.EXPECT().
			GetRevisionsV1(entity.GetRevisionsData{Claims: revClaims}, defPagination.IncrementLimit()).
			Return(expected, nil)

		result, err := uc.GetRevisionsV1(entity.GetRevisionsData{Claims: revClaims}, defPagination)

		require.NoError(t, err)
		require.Equal(t, expected, result.Data)
	})

	t.Run("Failed to get revisions", func(t *testing.T) {
		mockError := errors.New("mock error")
		accessorMock.EXPECT().CheckAccess(gomock.Any(), gomock.Any(), gomock.Any()).Return(accessor.AccessGranted, nil)
		repo.EXPECT().
			GetRevisionsV1(entity.GetRevisionsData{Claims: revClaims}, defPagination.IncrementLimit()).
			Return(nil, mockError)

		_, err := uc.GetRevisionsV1(entity.GetRevisionsData{Claims: revClaims}, defPagination)

		require.Error(t, err)
		require.Equal(t, mockError, errors.Unwrap(err))
	})
}
