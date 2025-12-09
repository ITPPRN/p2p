package auth

import (
	"bytes"
	"encoding/json"
	"net/http/httptest"
	"testing"

	"github.com/Nerzal/gocloak/v13"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	"p2p-back-end/configs"
	"p2p-back-end/logs"
	"p2p-back-end/modules/auth/controller"
	"p2p-back-end/modules/auth/repository"
	"p2p-back-end/modules/auth/service"
	"p2p-back-end/modules/entities/models"
)

func TestAuthIntegration(t *testing.T) {
	// ปิด log
	logs.Logger = zap.NewNop()

	// --------------------
	// setup database (SQLite in-memory)
	// --------------------
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	require.NoError(t, err)
	err = db.AutoMigrate(&models.UserEntity{})
	require.NoError(t, err)

	authRepo := repository.NewAuthRepositoryDB(db)

	// --------------------
	// setup Keycloak config จริง
	// --------------------
	cfg := &configs.Config{
		KeyCloak: configs.KeyCloak{
			Host:          "http://localhost", // เปลี่ยนเป็น host Keycloak จริง
			Port:          "9080",      // port Keycloak
			RealmName:     "P2P-service",
			ClientID:      "p2p-client",
			ClientSecret:  "4J30ZBRYRUTJcw8YlfUrTeA4qxKitBuu",
			AdminUsername: "admin",
			AdminPassword: "secure+admin+password+p2p",
		},
	}

	// --------------------
	// service ใช้ Keycloak จริง
	// --------------------
	keycloakClient := gocloak.NewClient(cfg.KeyCloak.Host + ":" + cfg.KeyCloak.Port)
	authSrv := service.NewAuthService(keycloakClient, cfg, authRepo)

	// --------------------
	// controller + fiber app
	// --------------------
	app := fiber.New()
	controller.NewUserController(app.Group("/v1/auth"), authSrv)

	// --------------------
	// Test Register
	// --------------------
	registerReq := models.RegisterReq{
		FirstName: "Integrat1ion",
		LastName:  "T1est",
		Email:     "int_te1st@example.com",
		Username:  "int_test11",
		Password:  "pass121134",
		Role:      "employee",
	}
	body, _ := json.Marshal(registerReq)
	req := httptest.NewRequest("POST", "/v1/auth/register", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	resp, _ := app.Test(req, -1)
	require.Equal(t, 200, resp.StatusCode)

	// --------------------
	// Test Login
	// --------------------
	loginReq := models.LoginReq{
		Username: "int_test11",
		Password: "pass121134",
	}
	bodyLogin, _ := json.Marshal(loginReq)
	reqLogin := httptest.NewRequest("POST", "/v1/auth/login", bytes.NewReader(bodyLogin))
	reqLogin.Header.Set("Content-Type", "application/json")
	respLogin, _ := app.Test(reqLogin, -1)
	require.Equal(t, 200, respLogin.StatusCode)

	// --------------------
	// Test repository directly
	// --------------------
	user := models.UserEntity{
		ID:        uuid.NewString(),
		Username:  registerReq.Username,
		Email:     registerReq.Email,
		FirstName: registerReq.FirstName,
		LastName:  registerReq.LastName,
		Role:      registerReq.Role,
		Password:  "hashed-pass",
	}
	db.Create(&user)

	exist, err := authRepo.IsUserExistByID(user.ID)
	require.NoError(t, err)
	require.True(t, exist)
}
