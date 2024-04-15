package usecase

import (
	"errors"
	"fmt"
	"single-window/internal/entity"
	"single-window/pkg/accessor"
)

type OrganizationUseCase struct {
	repo IOrganizationRepo
	acc  accessor.IAccessor
}

func NewOrganizationUseCase(repo IOrganizationRepo, acc accessor.IAccessor) *OrganizationUseCase {
	return &OrganizationUseCase{
		repo: repo,
		acc:  acc,
	}
}

func (ouc *OrganizationUseCase) GetOrganizationsV1(rev entity.GetOrganizationsData, pag entity.Pagination) (*entity.PaginatedOrganizations, error) {
	const (
		resource = "organizations"
		action   = "get"
	)
	accessStatus, err := ouc.acc.CheckAccess(rev.Claims.UserRole, resource, action)
	if err != nil {
		return nil, fmt.Errorf("GetOrganizationsV1 - UseCase %w", err)
	}

	if accessStatus != accessor.AccessGranted {
		return nil, errors.New("GetOrganizationsV1 - UseCase Access Denied")
	}

	organizations, err := ouc.repo.GetOrganizationsV1(rev, pag.IncrementLimit())
	if err != nil {
		return nil, fmt.Errorf("GetOrganizationsV1 - UseCase %w", err)
	}

	hasMoreData := len(organizations) > pag.Limit
	if hasMoreData {
		organizations = organizations[:len(organizations)-1]
	}

	return &entity.PaginatedOrganizations{
		Data:        organizations,
		HasMoreData: hasMoreData,
	}, nil
}
