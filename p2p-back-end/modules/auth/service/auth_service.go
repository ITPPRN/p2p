package service

import (
	"context"
	"fmt"
	"strings"

	"github.com/Nerzal/gocloak/v13"

	"p2p-back-end/configs"
	"p2p-back-end/logs"
	"p2p-back-end/modules/entities/models"
	"p2p-back-end/pkg/errs"
)

type authService struct {
	keycloak *gocloak.GoCloak
	cfg      *configs.Config
	authRepo models.AuthRepository
}

func NewAuthService(
	keycloak *gocloak.GoCloak,
	cfg *configs.Config,
	authRepo models.AuthRepository,
) models.AuthService {
	return &authService{keycloak, cfg, authRepo}
}

func (s authService) Register(req *models.RegisterReq) (string, error) {
	ctx := context.Background()

	token, err := s.keycloak.LoginAdmin(ctx, s.cfg.KeyCloak.AdminUsername, s.cfg.KeyCloak.AdminPassword, "master")
	if err != nil {
		logs.Error(err)
		return "", errs.NewUnexpectedError()
	}

	user := gocloak.User{
		FirstName: gocloak.StringP(req.FirstName),
		LastName:  gocloak.StringP(req.LastName),
		Email:     gocloak.StringP(req.Email),
		Enabled:   gocloak.BoolP(true),
		Username:  gocloak.StringP(req.Username),
	}

	roleName := strings.ToLower(req.Role)

	role, err := s.keycloak.GetRealmRole(ctx, token.AccessToken, s.cfg.KeyCloak.RealmName, roleName)
	// if err != nil {
	// 	logs.Error(err)
	// 	return "", errs.NewUnexpectedError()
	// }
	if err != nil {
		// ถ้าเกิด Error: 404 Not Found (Could not find role)
		// เราจะทำการสร้าง Role นั้นขึ้นมา
		if strings.Contains(err.Error(), "404 Not Found") || strings.Contains(err.Error(), "Could not find role") {
			logs.Warnf("Role '%s' not found in Keycloak. Attempting to create it automatically.", roleName)

			// สร้าง Role object ใหม่
			newRole := gocloak.Role{
				Name:        gocloak.StringP(roleName),
				Description: gocloak.StringP(fmt.Sprintf("Auto-generated role for user type: %s", roleName)),
			}

			// เรียก API สร้าง Role
			_, err = s.keycloak.CreateRealmRole(ctx, token.AccessToken, s.cfg.KeyCloak.RealmName, newRole)
			if err != nil {
				logs.Errorf("FATAL: Failed to auto-create role '%s': %v", roleName, err)
				return "", errs.NewUnexpectedError()
			}
			logs.Infof("✅ Successfully auto-created role: %s", roleName)

			// ดึง Role ที่สร้างใหม่มาอีกครั้งเพื่อใช้ในขั้นตอนต่อไป
			role, err = s.keycloak.GetRealmRole(ctx, token.AccessToken, s.cfg.KeyCloak.RealmName, roleName)
			if err != nil {
				logs.Errorf("Failed to retrieve newly created role '%s': %v", roleName, err)
				return "", errs.NewUnexpectedError()
			}

		} else {
			// ถ้าเป็น Error อื่นๆ ที่ไม่ใช่ 404 (เช่น Connection Error) ให้ Fail
			logs.Error(err)
			return "", errs.NewUnexpectedError()
		}
	}

	userID, err := s.keycloak.CreateUser(ctx, token.AccessToken, s.cfg.KeyCloak.RealmName, user)
	if err != nil {
		logs.Error(err)
		return "", errs.NewUnexpectedError()
	}

	err = s.keycloak.SetPassword(ctx, token.AccessToken, userID, s.cfg.KeyCloak.RealmName, req.Password, false)
	if err != nil {
		logs.Error(err)
		return "", errs.NewUnexpectedError()
	}

	// กำหนดเฉพาะบทบาทที่ระบุใน req.Role
	err = s.keycloak.AddRealmRoleToUser(ctx, token.AccessToken, s.cfg.KeyCloak.RealmName, userID, []gocloak.Role{*role})
	if err != nil {
		logs.Error(err)
		return "", errs.NewUnexpectedError()
	}

	// ลบบทบาท default-roles-master ถ้ามันถูกเพิ่มโดยอัตโนมัติ
	defaultRole, err := s.keycloak.GetRealmRole(ctx, token.AccessToken, s.cfg.KeyCloak.RealmName, "default-roles-"+s.cfg.KeyCloak.RealmName)
	if err == nil {
		err = s.keycloak.DeleteRealmRoleFromUser(ctx, token.AccessToken, s.cfg.KeyCloak.RealmName, userID, []gocloak.Role{*defaultRole})
		if err != nil {
			logs.Error(err)
			return "", errs.NewUnexpectedError()
		}
	}

	_, err = s.keycloak.Login(ctx, s.cfg.KeyCloak.ClientID, s.cfg.KeyCloak.ClientSecret, s.cfg.KeyCloak.RealmName, req.Username, req.Password)
	if err != nil {
		logs.Error(err)
		return "", errs.NewLoginFailedError()
	}
	return userID, nil
}

func (s authService) Login(req *models.LoginReq) (string, error) {

	ctx := context.Background()

	token, err := s.keycloak.Login(ctx, s.cfg.KeyCloak.ClientID, s.cfg.KeyCloak.ClientSecret, s.cfg.KeyCloak.RealmName, req.Username, req.Password)
	if err != nil {
		logs.Error(err)
		return "", errs.NewLoginFailedError()
	}

	return token.AccessToken, nil
}
