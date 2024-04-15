package usecase

import (
	"errors"
	"fmt"
	"single-window/internal/entity"
	"single-window/pkg/jwt"
	"time"
)

type authUseCase struct {
	w    IWbxAuthWebApi
	repo IAuthRepo
	g    jwt.IJWTGenerator
	ttl  time.Duration
}

func NewAuthUseCase(w IWbxAuthWebApi, repo IAuthRepo, g jwt.IJWTGenerator, ttl string) (IAuthUseCase, error) {
	tokenTTL, err := time.ParseDuration(ttl)
	if err != nil {
		return nil, fmt.Errorf("cannot parse duration: %w", err)
	}

	return &authUseCase{
		w:    w,
		repo: repo,
		g:    g,
		ttl:  tokenTTL,
	}, nil
}

func (ah *authUseCase) AuthHandle(req entity.AuthUserJSONBody, h entity.AuthUserParams) (*entity.TokenData, error) {
	res, err := ah.w.CheckAuthCode(req, h)
	if err != nil {
		return nil, fmt.Errorf("cannot check auth code: %w", err)
	}

	if res.Result != 0 {
		if res.Error != nil {
			return nil, fmt.Errorf("failed to check auth code: %s", *res.Error)
		}
		return nil, errors.New("failed to check auth code: unknown error")
	}

	claims, err := ah.repo.GetUserByPhone(req.PhoneNumber)
	if err != nil {
		return nil, fmt.Errorf("cannot get user by phone: %w", err)
	}

	token, err := ah.g.GenerateToken(claims, ah.ttl)
	if err != nil {
		return nil, fmt.Errorf("cannot generate token: %w", err)
	}

	return &entity.TokenData{
		Claims: *claims,
		Token:  *token,
		Ttl:    int(ah.ttl),
	}, nil
}
