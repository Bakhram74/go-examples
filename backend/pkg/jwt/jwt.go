package jwt

import (
	"crypto/ed25519"
	"fmt"
	"time"

	"github.com/go-jose/go-jose/v3"
	"github.com/go-jose/go-jose/v3/jwt"
)

type jwtParser struct {
	publicKey *ed25519.PublicKey
}

func NewParser(publicKey []byte) (IJWTParser, error) {
	pub, err := ParseEd25519PublicKeyFromPEM(publicKey)
	if err != nil {
		return nil, err
	}

	return &jwtParser{publicKey: pub}, nil
}

func (p *jwtParser) ParseToken(accessToken string, claims any) error {
	token, err := jwt.ParseSigned(accessToken)
	if err != nil {
		return fmt.Errorf("cannot parse token: %w", err)
	}
	requiredClaims := jwt.Claims{}
	if err := token.Claims(*p.publicKey, &requiredClaims, claims); err != nil {
		return fmt.Errorf("cannot get claims: %w", err)
	}
	if err := requiredClaims.Validate(jwt.Expected{
		Time:     time.Now(),
		Issuer:   "",
		Subject:  "",
		Audience: jwt.Audience{},
		ID:       "",
	}); err != nil {
		return fmt.Errorf("error validating jwt claims: %w", err)
	}
	return nil
}

type jwtGenerator struct {
	privateKey *ed25519.PrivateKey
}

func NewGenerator(privateKey []byte) (IJWTGenerator, error) {
	pri, err := ParseEd25519PrivateKeyFromPEM(privateKey)
	if err != nil {
		return nil, err
	}

	return &jwtGenerator{privateKey: pri}, nil
}

func (p *jwtGenerator) GenerateToken(claims any, ttl time.Duration) (*string, error) {
	signer, err := jose.NewSigner(jose.SigningKey{Algorithm: jose.EdDSA, Key: *p.privateKey}, nil)
	if err != nil {
		return nil, fmt.Errorf("can't create Signer for token: %w", err)
	}

	requiredClaims := jwt.Claims{
		IssuedAt: jwt.NewNumericDate(time.Now()),
		Expiry:   jwt.NewNumericDate(time.Now().Add(ttl)),
	}

	token, err := jwt.Signed(signer).
		Claims(requiredClaims).
		Claims(claims).
		CompactSerialize()
	if err != nil {
		return nil, fmt.Errorf("can't create token: %w", err)
	}
	return &token, nil
}
