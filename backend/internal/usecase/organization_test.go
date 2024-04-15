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

func initOrganizationsMocks(t *testing.T) (*usecase.OrganizationUseCase, *MockIOrganizationRepo, *accessortest.MockIAccessor) {
	t.Helper()

	mockCtl := gomock.NewController(t)
	defer mockCtl.Finish()

	mockRepo := NewMockIOrganizationRepo(mockCtl)
	accessorMock := accessortest.NewMockIAccessor(mockCtl)
	mockUseCase := usecase.NewOrganizationUseCase(mockRepo, accessorMock)

	return mockUseCase, mockRepo, accessorMock
}

var defOrgClaims = entity.Claims{
	UserId:         "11111111-AAAA-AAAA-AAAA-AAAAAAAAAAAA",
	OrganizationId: "11111111-AAAA-AAAA-AAAA-AAAAAAAAAAAA",
}

var defOrgPagination = entity.Pagination{
	Limit:  3,
	Offset: 0,
}

func TestGetOrganizationsV1(t *testing.T) {
	t.Parallel()

	uc, repo, accessorMock := initOrganizationsMocks(t)

	t.Run("Empty result", func(t *testing.T) {
		accessorMock.EXPECT().CheckAccess(gomock.Any(), gomock.Any(), gomock.Any()).Return(accessor.AccessGranted, nil)
		repo.EXPECT().
			GetOrganizationsV1(entity.GetOrganizationsData{Claims: defOrgClaims}, defOrgPagination.IncrementLimit()).
			Return(nil, nil)

		result, err := uc.GetOrganizationsV1(entity.GetOrganizationsData{Claims: defOrgClaims}, defOrgPagination)

		require.NoError(t, err)
		require.Nil(t, result.Data)
	})

	t.Run("Get organizations", func(t *testing.T) {
		expected := []entity.Organization{{}, {}, {}}
		accessorMock.EXPECT().CheckAccess(gomock.Any(), gomock.Any(), gomock.Any()).Return(accessor.AccessGranted, nil)
		repo.EXPECT().
			GetOrganizationsV1(entity.GetOrganizationsData{Claims: defOrgClaims}, defOrgPagination.IncrementLimit()).
			Return(expected, nil)

		result, err := uc.GetOrganizationsV1(entity.GetOrganizationsData{Claims: defOrgClaims}, defOrgPagination)

		require.NoError(t, err)
		require.Equal(t, expected, result.Data)
	})

	t.Run("Failed to get organizations", func(t *testing.T) {
		mockError := errors.New("mock error")
		accessorMock.EXPECT().CheckAccess(gomock.Any(), gomock.Any(), gomock.Any()).Return(accessor.AccessGranted, nil)
		repo.EXPECT().
			GetOrganizationsV1(entity.GetOrganizationsData{Claims: defOrgClaims}, defOrgPagination.IncrementLimit()).
			Return(nil, mockError)

		_, err := uc.GetOrganizationsV1(entity.GetOrganizationsData{Claims: defOrgClaims}, defOrgPagination)

		require.Error(t, err)
		require.Equal(t, mockError, errors.Unwrap(err))
	})
}
