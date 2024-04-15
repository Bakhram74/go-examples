package usecase

import (
	"single-window/internal/entity"
)

type (
	IDisputeUseCase interface {
		GetDisputesV1(claims entity.Claims, limits, offset int, status entity.EntityStatus) ([]entity.DisputeList, *bool, error)
		GetDisputeV1(disputeId *string, goodsId *int) (*entity.Dispute, error)
		CloseDisputeV1(claims *entity.Claims, disputeId string, guiltyWorkerIds []string) (*entity.CloseDispute, error)
		CreateDisputeV1(claims *entity.Claims, createDisputeReqBody *entity.CreateDisputeV1MultipartBody) (*entity.CreateDispute, error)
		GetShortagesV1(userId string, limits, offset int) ([]entity.Shortage, *bool, error)
	}

	IDisputeRepo interface {
		GetDisputesV1(claims entity.Claims, limits, offset int, status entity.EntityStatus) ([]entity.DisputeList, error)
		GetDisputeV1(disputeId *string, goodsId *int) (*entity.Dispute, error)
		CreateDisputeV1(createDisputeReqBody *entity.CreateDisputeV1MultipartBody, userId string, disputeAttachmentPath *string) (*entity.CreateDispute, error)
		DeleteDisputeV1(dispute *entity.CreateDispute) (*entity.CreateDispute, error)
		GetGuiltyWorkersV1(disputeId string) ([]string, error)
		CloseDisputeV1(disputeId string, userId []string) (*entity.CloseDispute, error)
		GetShortagesV1(userId string, limits, offset int) ([]entity.Shortage, error)
	}

	IMessageUseCase interface {
		GetMessagesV1(claims *entity.Claims, disputeId string, limits, offset int) ([]entity.Message, *bool, error)
		CreateMessageV1(claims *entity.Claims, createMsgReqBody entity.CreateMessageV1MultipartBody) (*entity.CreateMsgResponse, error)
	}

	IMessageRepo interface {
		GetMessagesV1(userId, disputeId string, limits, offset int) ([]entity.Message, error)
		CreateMessageV1(createMsgReqBody entity.CreateMessageV1MultipartBody, userId string, attachPath *string) (*entity.CreateMsgResponse, error)
	}

	IAuthUseCase interface {
		AuthHandle(req entity.AuthUserJSONBody, h entity.AuthUserParams) (*entity.TokenData, error)
	}

	IAuthRepo interface {
		GetUserByPhone(phone string) (*entity.Claims, error)
	}

	IWbxAuthWebApi interface {
		CheckAuthCode(body entity.AuthUserJSONBody, h entity.AuthUserParams) (*entity.WbxAuthCodeCheckResponse, error)
	}

	IRevisionUseCase interface {
		GetRevisionsV1(rev entity.GetRevisionsData, pag entity.Pagination) (*entity.PaginatedRevisions, error)
		CreateRevisionV1(claims *entity.Claims, body *entity.CreateRevisionV1MultipartBody) (*entity.Revision, error)
		GetCorrespondencesV1(rev entity.GetCorrespondencesData, pag entity.Pagination) ([]entity.Correspondence, *bool, error)
		GetRevisionV1(claims *entity.Claims, revisionId string) (*entity.Revision, error)
	}

	IRevisionRepo interface {
		GetRevisionsV1(rev entity.GetRevisionsData, pag entity.Pagination) ([]entity.Revision, error)
		CreateRevisionV1(body *entity.CreateRevisionV1MultipartBody, userId string, revisionAttachmentPath *string) (*entity.Revision, error)
		GetCorrespondencesV1(revisionId string, pag entity.Pagination) ([]entity.Correspondence, error)
		GetRevisionV1(revisionId string) (*entity.Revision, error)
	}

	IOrganizationUseCase interface {
		GetOrganizationsV1(rev entity.GetOrganizationsData, pag entity.Pagination) (*entity.PaginatedOrganizations, error)
	}

	IOrganizationRepo interface {
		GetOrganizationsV1(rev entity.GetOrganizationsData, pag entity.Pagination) ([]entity.Organization, error)
	}
)
