package repo

import (
	"fmt"
	"single-window/internal/entity"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

const (
	revisionsTable       = "disputes.revision"
	correspondencesTable = "disputes.correspondence"
)

type RevisionRepo struct {
	db *gorm.DB
}

func NewRevisionRepo(db *gorm.DB) *RevisionRepo {
	return &RevisionRepo{db}
}

func (r *RevisionRepo) GetRevisionsV1(rev entity.GetRevisionsData, pag entity.Pagination) ([]entity.Revision, error) {
	var revisions []entity.Revision

	db := r.db.Table("disputes.revision AS r").
		Select(
			"r.*",
			"u.fullname AS worker_name",
			"o.organization_title AS organization_title").
		Joins("LEFT JOIN subjects.user AS u ON r.worker_id = u.user_id").
		Joins("LEFT JOIN subjects.organization AS o ON o.organization_id = r.organization_id")

	if rev.DisputeID != nil {
		db = db.
			Joins("LEFT JOIN disputes.dispute_role dr ON r.dispute_id = dr.dispute_id").
			Where("r.dispute_id = ?", rev.DisputeID).
			Where("dr.user_id = ?", rev.Claims.UserId).
			Where("dr.dispute_role IN ('responsible_person', 'arbitr')")
	} else {
		db = db.Where("r.organization_id = ?", rev.Claims.OrganizationId)

		switch rev.Status {
		case entity.Any:
			db = db.Where(
				r.db.Where("r.status = ?", entity.Opened).Or(
					"r.worker_id = ? AND r.status IN ?",
					rev.Claims.UserId,
					[]entity.EntityStatus{entity.InWork, entity.Closed},
				),
			)
		case entity.Opened:
			db = db.Where("r.status = ?", rev.Status)
		case entity.Closed, entity.InWork:
			db = db.Where("r.status = ?", rev.Status).Where("r.worker_id = ?", rev.Claims.UserId)
		}
	}

	err := db.Limit(pag.Limit).
		Offset(pag.Offset).
		Order("r.created_at DESC").
		Find(&revisions).Error
	if err != nil {
		return nil, fmt.Errorf("failed to get revisions, err: %w", err)
	}

	return revisions, nil
}

func (r *RevisionRepo) CreateRevisionV1(body *entity.CreateRevisionV1MultipartBody, userId string, revisionAttachmentPath *string) (*entity.Revision, error) {
	createdAt := time.Now()
	revisionId := uuid.NewString()
	revisionToCreate := entity.Revision{
		CreatedAt:      createdAt,
		DisputeId:      body.DisputeId,
		RevisionId:     revisionId,
		OrganizationId: body.OrganizationId,
		Status:         entity.Opened,
		WorkerId:       nil,
		InWorkAt:       nil,
		ClosedAt:       nil,
	}

	correspondenceToCreate := entity.Correspondence{
		AttachmentPath:   revisionAttachmentPath,
		CorrespondenceId: uuid.NewString(),
		CreatedAt:        createdAt,
		MessageBody:      &body.MessageBody,
		RevisionId:       revisionId,
		SenderId:         userId,
	}

	err := r.db.Transaction(func(tx *gorm.DB) error {
		err := tx.Table(revisionsTable).Omit("WorkerName", "OrganizationTitle").Create(revisionToCreate).Error
		if err != nil {
			return fmt.Errorf("failed to create revision, error: %w", err)
		}

		err = tx.Table(correspondencesTable).Create(correspondenceToCreate).Error
		if err != nil {
			return fmt.Errorf("failed to create correspondence, error: %w", err)
		}

		return nil
	})

	return &revisionToCreate, err
}

func (r *RevisionRepo) GetCorrespondencesV1(revisionId string, pag entity.Pagination) ([]entity.Correspondence, error) {
	var correspondences []entity.Correspondence

	err := r.db.Table("disputes.correspondence AS c").
		Select("c.*").
		Where("c.revision_id = ?", revisionId).
		Limit(pag.Limit).
		Offset(pag.Offset).
		Order("c.created_at DESC").
		Find(&correspondences).
		Error
	if err != nil {
		return nil, fmt.Errorf("failed to get correspondences, err: %w", err)
	}

	return correspondences, nil
}

func (r *RevisionRepo) GetRevisionV1(revisionId string) (*entity.Revision, error) {
	var revision entity.Revision

	db := r.db.Table("disputes.revision AS r").
		Select(
			"r.*",
			"u.fullname AS worker_name",
			"o.organization_title AS organization_title").
		Joins("LEFT JOIN subjects.user AS u ON r.worker_id = u.user_id").
		Joins("LEFT JOIN subjects.organization AS o ON o.organization_id = r.organization_id").
		Where("r.revision_id = ?", revisionId)

	err := db.First(&revision).Error
	if err != nil {
		return nil, fmt.Errorf("failed to get revision, err: %w", err)
	}

	return &revision, err
}
