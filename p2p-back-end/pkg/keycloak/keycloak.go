package keycloak

import (
	"p2p-back-end/configs"
	"p2p-back-end/logs"
	"p2p-back-end/pkg/utils"

	"github.com/Nerzal/gocloak/v13"
	"go.uber.org/zap"
)

func NewKeyCloakClient(cfg *configs.Config) *gocloak.GoCloak {

	url, err := utils.UrlBuilder("keycloak", cfg)
	if err != nil {
		logs.Error(zap.Error(err))
	}
	return gocloak.NewClient(url)
}
