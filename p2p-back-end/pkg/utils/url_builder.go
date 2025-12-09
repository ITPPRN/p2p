package utils

import (
	"errors"
	"fmt"

	"p2p-back-end/configs"
)

func UrlBuilder(urlType string, cfg *configs.Config) (string, error) {

	var url string

	switch urlType {
	case "fiber":
		url = fmt.Sprintf(":%s", cfg.App.Port)
	case "postgres":
		url = fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s  sslmode=disable TimeZone=Asia/Bangkok",
			cfg.Postgres.Host,
			cfg.Postgres.Port,
			cfg.Postgres.Username,
			cfg.Postgres.Password,
			cfg.Postgres.DatabaseName,
			cfg.Postgres.SslMode,
		)
	case "redis":
		url = fmt.Sprintf("%s:%s", cfg.Redis.Host, cfg.Redis.Port)
	case "keycloak":
		url = fmt.Sprintf("http://%s:%s", cfg.KeyCloak.Host, cfg.KeyCloak.Port)
		// url = fmt.Sprintf("%s:%s", cfg.KeyCloak.Host, cfg.KeyCloak.Port)
		//url = "http://localhost:8080"
	default:
		err := fmt.Sprintf("error,url builder Unknown url type: %s", urlType)
		return "", errors.New(err)
	}
	return url, nil
}
