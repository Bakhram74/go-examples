package jwt

import "time"

type (
	IJWTGenerator interface {
		GenerateToken(claims any, ttl time.Duration) (*string, error)
	}

	IJWTParser interface {
		ParseToken(accessToken string, claims any) error
	}
)
