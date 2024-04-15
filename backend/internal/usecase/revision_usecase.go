package usecase

import (
	"errors"
	"fmt"
	"single-window/internal/entity"
	"single-window/pkg/accessor"
	"single-window/pkg/s3"
	"single-window/pkg/utils"
)

type RevisionUseCase struct {
	repo     IRevisionRepo
	acc      accessor.IAccessor
	s3Client s3.IS3Client
}

func NewRevisionUseCase(repo IRevisionRepo, acc accessor.IAccessor, s3Client s3.IS3Client) *RevisionUseCase {
	return &RevisionUseCase{
		repo:     repo,
		acc:      acc,
		s3Client: s3Client,
	}
}

func (ruc *RevisionUseCase) GetRevisionsV1(rev entity.GetRevisionsData, pag entity.Pagination) (*entity.PaginatedRevisions, error) {
	const (
		resource = "revisions"
		action   = "get"
	)
	accessStatus, err := ruc.acc.CheckAccess(rev.Claims.UserRole, resource, action)
	if err != nil {
		return nil, fmt.Errorf("GetRevisionsV1 - UseCase %w", err)
	}

	if accessStatus != accessor.AccessGranted {
		return nil, errors.New("GetRevisionsV1 - UseCase Access Denied")
	}

	revisions, err := ruc.repo.GetRevisionsV1(rev, pag.IncrementLimit())
	if err != nil {
		return nil, fmt.Errorf("GetRevisionsV1 - UseCase %w", err)
	}

	hasMoreData := len(revisions) > pag.Limit
	if hasMoreData {
		revisions = revisions[:len(revisions)-1]
	}

	return &entity.PaginatedRevisions{
		Data:        revisions,
		HasMoreData: hasMoreData,
	}, nil
}

func (ruc *RevisionUseCase) CreateRevisionV1(claims *entity.Claims, body *entity.CreateRevisionV1MultipartBody) (*entity.Revision, error) {
	// TODO: filePath лучше составлять как correspondence_id и filename
	attachmentPath, err := utils.UploadFile(ruc.s3Client, claims.UserId, body.File)
	if err != nil {
		return nil, fmt.Errorf("failed to put file in bucket, err: %w", err)
	}

	createdRevision, err := ruc.repo.CreateRevisionV1(body, claims.UserId, attachmentPath)
	if err != nil {
		return nil, fmt.Errorf("failed to create revision, err: %w", err)
	}

	return createdRevision, nil
}

func (ruc *RevisionUseCase) GetCorrespondencesV1(rev entity.GetCorrespondencesData, pag entity.Pagination) ([]entity.Correspondence, *bool, error) {
	accessStatus, err := ruc.acc.CheckAccess(rev.Claims.UserRole, "correspondences", "get")
	if err != nil {
		return nil, nil, fmt.Errorf("GetCorrespondencesV1 - UseCase %w", err)
	}

	if accessStatus != accessor.AccessGranted {
		return nil, nil, errors.New("GetCorrespondencesV1 - UseCase Access Denied")
	}

	correspondences, err := ruc.repo.GetCorrespondencesV1(rev.RevisionID, pag.IncrementLimit())
	if err != nil {
		return nil, nil, fmt.Errorf("GetCorrespondencesV1 - UseCase %w", err)
	}

	hasMoreData := len(correspondences) > pag.Limit
	if hasMoreData {
		correspondences = correspondences[:len(correspondences)-1]
	}

	return correspondences, &hasMoreData, err
}

func (ruc *RevisionUseCase) GetRevisionV1(claims *entity.Claims, revesionId string) (*entity.Revision, error) {
	accessStatus, err := ruc.acc.CheckAccess(claims.UserRole, "revisions", "get")
	if err != nil {
		return nil, fmt.Errorf("GetRevisionV1 - UseCase %w", err)
	}

	if accessStatus != accessor.AccessGranted {
		return nil, errors.New("GetRevisionV1 - UseCase Access Denied")
	}

	revision, err := ruc.repo.GetRevisionV1(revesionId)
	if err != nil {
		return nil, fmt.Errorf("failed to get revision, err: %w", err)
	}

	return revision, nil
}
