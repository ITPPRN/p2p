package servers

import (
	"github.com/Nerzal/gocloak/v13"
	"github.com/go-redis/redis/v8"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"

	"p2p-back-end/configs"
	"p2p-back-end/logs"
	"p2p-back-end/pkg/utils"
)

type server struct {
	App      *fiber.App
	Db       *gorm.DB
	Cfg      *configs.Config
	Redis    *redis.Client
	Keycloak *gocloak.GoCloak
}

func NewServer(
	cfg *configs.Config,
	db *gorm.DB,
	redis *redis.Client,
	keycloak *gocloak.GoCloak,

) *server {
	return &server{
		App: fiber.New( /*fiber.Config{
			Prefork: cfg.App.Mode == "release", // production Prefork : true
			// Prefork: false,
		}*/),
		Db:       db,
		Cfg:      cfg,
		Redis:    redis,
		Keycloak: keycloak,
	}
}

func (s *server) Start() {

	if err := s.Handlers(); err != nil {
		logs.Error(err)
		panic(err.Error())
	}

	fiberConnURL, err := utils.UrlBuilder("fiber", s.Cfg)
	if err != nil {
		logs.Error(err)
		panic(err.Error())
	}

	port := s.Cfg.App.Port
	mode := s.Cfg.App.Mode

	if mode == "release" {
		logs.Info("Fiber server start in production mode at port " + port)
		if err := s.App.Listen(fiberConnURL); err != nil {
			logs.Error(err)
			// ตามความเหมาะสม, คุณอาจต้องการที่จะส่งคืนหรือปิดแอปให้มีความสุขตามต้องการของแอป
		}
	} else {
		logs.Info("Fiber server start in debug mode at port " + port)
		if err := s.App.Listen(fiberConnURL); err != nil {
			logs.Error(err)
			// ตามความเหมาะสม, คุณอาจต้องการที่จะส่งคืนหรือปิดแอปให้มีความสุขตามต้องการของแอป
		}
	}

}
