package repo

import (
	"fmt"
	"single-window/internal/entity"

	"gorm.io/gorm"
)

type OrganizationRepo struct {
	db *gorm.DB
}

func NewOrganizationRepo(db *gorm.DB) *OrganizationRepo {
	return &OrganizationRepo{db}
}

func (r *OrganizationRepo) GetOrganizationsV1(rev entity.GetOrganizationsData, pag entity.Pagination) ([]entity.Organization, error) {
	var organizations []entity.Organization

	db := r.db.Table("subjects.organization AS o").
		Select("o.*").
		Where("o.is_del = ?", false)
	if rev.SearchToken != nil {
		db = db.Where("o.organization_title ILIKE ?", fmt.Sprintf("%%%s%%", *rev.SearchToken))
	}
	if rev.OrganizationCode != nil {
		db = db.Where("o.organization_code = ?", *rev.OrganizationCode)
	}

	err := db.Limit(pag.Limit).
		Offset(pag.Offset).
		Order("o.organization_title ASC").
		Find(&organizations).Error
	if err != nil {
		return nil, fmt.Errorf("failed to get organizations, err: %w", err)
	}

	return organizations, nil
}
