package config

import (
	"fmt"
	"os"
	"single-window/pkg/httpclient/fasthttpclient"
	"strings"

	"github.com/joho/godotenv"
	"github.com/knadh/koanf/parsers/yaml"
	"github.com/knadh/koanf/providers/env"
	"github.com/knadh/koanf/providers/file"
	"github.com/knadh/koanf/v2"
	"github.com/pkg/errors"
)

const pathToMainFolder = "../.."

var (
	pathToConfig = fmt.Sprintf("%s/config/config.yml", pathToMainFolder)
	pathToEnv    = fmt.Sprintf("%s/.env", pathToMainFolder)
)

type (
	Config struct {
		App
		HTTP
		Log
		PG
		MINIO
		JWT
		AuthClient
		AuthCookie
		Casbin
	}

	App struct {
		Name    string
		Version string
		Mode    string
	}

	HTTP struct {
		Port string
	}

	Log struct {
		Level string
	}

	PG struct {
		PoolMax int
		URL     string
	}

	MINIO struct {
		URL             string
		AccessKeyID     string
		SecretAccessKey string
		UseSSL          bool
	}

	JWT struct {
		TokenTTL   string
		PrivateKey string
		PublicKey  string
	}

	AuthCookie struct {
		Name   string
		Path   string
		Domain string
	}

	AuthClient struct {
		WBURL string
		fasthttpclient.FastHttpConfig
	}

	Casbin struct {
		ConfigPath string
	}
)

func NewConfig() (*Config, error) {
	isLoadedFromDotEnv, err := loadDotEnv()
	if err != nil {
		return nil, fmt.Errorf("failed to read env file, err: %w", err)
	}

	k := koanf.New(".")

	if err := k.Load(file.Provider(pathToConfig), yaml.Parser()); err != nil {
		return nil, fmt.Errorf("failed to load config, err: %w", err)
	}

	providerCb := func(s string) string {
		return strings.ReplaceAll(strings.ToLower(s), "_", ".")
	}

	if err := k.Load(env.Provider("", ".", providerCb), nil); err != nil {
		return nil, fmt.Errorf("failed to load config parser, err: %w", err)
	}

	var conf Config
	if err := k.Unmarshal("", &conf); err != nil {
		return nil, fmt.Errorf("failed to decode config into struct, err: %w", err)
	}

	if conf.App.Mode != "local" && isLoadedFromDotEnv {
		return nil, errors.New(`app with an ".env" file can only be launched in "local" mode`)
	}

	return &conf, nil
}

func loadDotEnv() (bool, error) {
	if _, err := os.Stat(pathToEnv); os.IsNotExist(err) {
		return false, nil
	}
	if err := godotenv.Load(pathToEnv); err != nil {
		return true, fmt.Errorf("cannot load env err: %w", err)
	}

	return true, nil
}
