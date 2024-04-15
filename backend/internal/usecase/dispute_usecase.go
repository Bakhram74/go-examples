package usecase

import (
	"bytes"
	"errors"
	"fmt"
	"single-window/internal/entity"
	"single-window/pkg/accessor"
	"single-window/pkg/s3"

	"github.com/minio/minio-go/v7"
)

type DisputeUseCase struct {
	repo     IDisputeRepo
	s3Client s3.IS3Client
	acc      accessor.IAccessor
}

func NewDisputeUseCase(dr IDisputeRepo, s3Client s3.IS3Client, acc accessor.IAccessor) *DisputeUseCase {
	return &DisputeUseCase{
		repo:     dr,
		s3Client: s3Client,
		acc:      acc,
	}
}

func (duc *DisputeUseCase) GetDisputesV1(claims entity.Claims, limits, offset int, status entity.EntityStatus) ([]entity.DisputeList, *bool, error) {
	const (
		resource = "disputes"
		action   = "get"
	)
	accessStatus, err := duc.acc.CheckAccess(claims.UserRole, resource, action)
	if err != nil {
		return nil, nil, fmt.Errorf("GetDisputesV1 - UseCase %w", err)
	}

	if accessStatus != accessor.AccessGranted {
		return nil, nil, errors.New("GetDisputesV1 - UseCase Access Denied")
	}

	disputes, err := duc.repo.GetDisputesV1(claims, limits+1, offset, status)
	if err != nil {
		return nil, nil, fmt.Errorf("GetDisputesV1 - UseCase %w", err)
	}

	hasMoreData := len(disputes) > limits
	if hasMoreData {
		disputes = disputes[:len(disputes)-1]
	}

	return disputes, &hasMoreData, nil
}

func (duc *DisputeUseCase) GetDisputeV1(disputeId *string, goodsId *int) (*entity.Dispute, error) {
	if disputeId == nil && goodsId == nil {
		return nil, fmt.Errorf("param dispute_id or goods_id cannot be null")
	}

	dispute, err := duc.repo.GetDisputeV1(disputeId, goodsId)
	if err != nil {
		return nil, fmt.Errorf("GetDisputeV1 - UseCase - %w", err)
	}

	dispute.GuiltyWorkerNames, err = duc.repo.GetGuiltyWorkersV1(dispute.DisputeId)
	if err != nil {
		return dispute, fmt.Errorf("failed to get gulityWorkers - %w", err)
	}

	return dispute, nil
}

func (duc *DisputeUseCase) CreateDisputeV1(claims *entity.Claims, createDisputeReqBody *entity.CreateDisputeV1MultipartBody) (*entity.CreateDispute, error) {
	var disputeAttachmentPath *string
	if createDisputeReqBody.File != nil {
		fBytes, err := createDisputeReqBody.File.Bytes()
		if err != nil {
			return nil, fmt.Errorf("failed to get file bytes, err: %w", err)
		}

		filePath := fmt.Sprintf("%s/%s", claims.UserId, createDisputeReqBody.File.Filename())
		// TODO: перед этим, нужно проверить, нет ли там уже файла с таким же названием
		// TODO: настроить expire date для файлов, подумать над тем, как будут удаляться файлы из хранилища
		uploadInfo, err := duc.s3Client.PutObjectToBucket(filePath, bytes.NewReader(fBytes), int64(len(fBytes)), minio.PutObjectOptions{})
		if err != nil {
			return nil, fmt.Errorf("failed to put file in bucket, err: %w", err)
		}
		disputeAttachmentPath = &uploadInfo.Key
	}

	createdDispute, err := duc.repo.CreateDisputeV1(createDisputeReqBody, claims.UserId, disputeAttachmentPath)
	if err != nil {
		return nil, fmt.Errorf("failed to create dispute, err: %w", err)
	}

	return createdDispute, nil
}

func (duc *DisputeUseCase) CloseDisputeV1(claims *entity.Claims, disputeId string, guiltyWorkerIds []string) (*entity.CloseDispute, error) {
	const (
		resource = "disputes"
		action   = "close"
	)
	accessStatus, err := duc.acc.CheckAccess(claims.UserRole, resource, action)
	if err != nil {
		return nil, fmt.Errorf("GetDisputesV1 - UseCase %w", err)
	}

	if accessStatus != accessor.AccessGranted {
		return nil, errors.New("GetDisputesV1 - UseCase Access Denied")
	}

	if claims.UserRole == "arbitr" {
		if disputeId == "" {
			return nil, fmt.Errorf("CloseDisputeV1 - UseCase - error disputeId has null")
		}
	} else if disputeId == "" || len(guiltyWorkerIds) == 0 {
		return nil, fmt.Errorf("CloseDisputeV1 - UseCase - error disputeId has null or guiltyWorkersIds empty")
	}

	dispute, err := duc.repo.CloseDisputeV1(disputeId, guiltyWorkerIds)
	if err != nil {
		return nil, fmt.Errorf("CloseDisputeV1 - UseCase - CloseDispute - %w", err)
	}

	return dispute, nil
}

func (duc *DisputeUseCase) GetShortagesV1(userId string, limits, offset int) ([]entity.Shortage, *bool, error) {
	shortages, err := duc.repo.GetShortagesV1(userId, limits+1, offset)
	if err != nil {
		return nil, nil, fmt.Errorf("GetShortagesV1 - UseCase - %w", err)
	}

	hasMoreData := len(shortages) > limits
	if hasMoreData {
		shortages = shortages[:len(shortages)-1]
	}

	return shortages, &hasMoreData, nil
}
