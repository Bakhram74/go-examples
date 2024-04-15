package jwt

import (
	"crypto/ed25519"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"fmt"
)

var (
	ErrKeyMustBePEMEncoded  = errors.New("invalid Key: Key must be PEM encoded PKCS1 or PKCS8 private key")
	ErrNotEd25519PrivateKey = errors.New("key is not a valid Ed25519 private key")
	ErrNotEd25519PublicKey  = errors.New("key is not a valid Ed25519 public key")
)

func ParseEd25519PrivateKeyFromPEM(key []byte) (*ed25519.PrivateKey, error) {
	var (
		err   error
		block *pem.Block
	)
	if block, _ = pem.Decode(key); block == nil {
		return nil, ErrKeyMustBePEMEncoded
	}

	var parsedKey interface{}
	if parsedKey, err = x509.ParsePKCS8PrivateKey(block.Bytes); err != nil {
		return nil, fmt.Errorf("error parsing private key: %w", err)
	}

	var pkey ed25519.PrivateKey
	var ok bool
	if pkey, ok = parsedKey.(ed25519.PrivateKey); !ok {
		return nil, ErrNotEd25519PrivateKey
	}

	return &pkey, nil
}

type parserFunc func(der []byte) (any, error)

func parseCertificate(der []byte) (any, error) {
	cert, err := x509.ParseCertificate(der)
	if err != nil {
		return nil, fmt.Errorf("cannot parsing cert: %w", err)
	}
	return cert.PublicKey, nil
}

func ParseEd25519PublicKeyFromPEM(key []byte) (*ed25519.PublicKey, error) {
	var (
		err   error
		block *pem.Block
	)
	if block, _ = pem.Decode(key); block == nil {
		return nil, ErrKeyMustBePEMEncoded
	}

	parserFn := []parserFunc{x509.ParsePKIXPublicKey, parseCertificate}

	var parsedKey interface{}
	for i := range parserFn {
		parsedKey, err = parserFn[i](block.Bytes)
		if err == nil {
			break
		}
	}
	if err != nil {
		return nil, fmt.Errorf("error parsing public key: %w", err)
	}

	var pkey ed25519.PublicKey
	var ok bool
	if pkey, ok = parsedKey.(ed25519.PublicKey); !ok {
		return nil, ErrNotEd25519PublicKey
	}

	return &pkey, nil
}
