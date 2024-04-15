package v1

import (
	"fmt"
	"single-window/config"
	"single-window/internal/usecase"
	"single-window/internal/usecase/repo"
	"single-window/internal/usecase/webapi"
	"single-window/pkg/accessor"
	"single-window/pkg/jwt"
	"single-window/pkg/logger"
	"single-window/pkg/s3"
	"single-window/pkg/swcontext"

	"gorm.io/gorm"
)

type routes struct {
	disputeRoutes
	authRoutes
	messageRoutes
	revisionRoutes
	organizationRoutes
}

func NewRoutes(
	db *gorm.DB, s3Client s3.IS3Client,
	cfg *config.Config, l logger.ILogger, swCtxAdapter swcontext.SwGinContext, acc accessor.IAccessor,
) (ServerInterface, error) {
	disputeRoutes := NewDisputeRoutes(
		usecase.NewDisputeUseCase(repo.NewDisputeRepo(db), s3Client, acc),
		l,
		swCtxAdapter,
	)

	messageRoutes := NewMessageRoutes(
		usecase.NewMessageUseCase(repo.NewMessageRepo(db), s3Client, acc),
		l,
		swCtxAdapter,
	)

	generator, err := jwt.NewGenerator([]byte(cfg.JWT.PrivateKey))
	if err != nil {
		return nil, fmt.Errorf("failed to create jwt generator, err: %w", err)
	}

	cli, err := webapi.NewWbAuthClient(cfg.AuthClient)
	if err != nil {
		return nil, fmt.Errorf("failed to create auth client, err: %w", err)
	}

	auc, err := usecase.NewAuthUseCase(cli,
		repo.NewAuthRepo(db), generator, cfg.JWT.TokenTTL)
	if err != nil {
		return nil, fmt.Errorf("failed to create auth usecase, err: %w", err)
	}

	authRoutes := NewAuthRoutes(auc, cfg.AuthCookie.Name, cfg.AuthCookie.Path, cfg.AuthCookie.Domain)

	revisionRoutes := newRevisionRoutes(
		usecase.NewRevisionUseCase(repo.NewRevisionRepo(db), acc, s3Client),
		l,
		swCtxAdapter,
	)

	organizationRoutes := newOrganizationRoutes(
		usecase.NewOrganizationUseCase(repo.NewOrganizationRepo(db), acc),
		l, swCtxAdapter,
	)

	return &routes{
		disputeRoutes:      *disputeRoutes,
		authRoutes:         *authRoutes,
		messageRoutes:      *messageRoutes,
		revisionRoutes:     *revisionRoutes,
		organizationRoutes: *organizationRoutes,
	}, nil
}
