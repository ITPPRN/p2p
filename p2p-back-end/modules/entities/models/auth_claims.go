package models

import "github.com/golang-jwt/jwt/v5"

// --- Structs (Unchanged) ---

type JWTPayload struct {
	Azp         string      `json:"azp"`
	Email       string      `json:"email"`
	Exp         int64       `json:"exp"`
	Iat         int64       `json:"iat"`
	ID          string      `json:"id"`
	Iss         string      `json:"iss"`
	Jti         string      `json:"jti"`
	Name        string      `json:"name"`
	RealmAccess RealmAccess `json:"realm_access"`
	Scope       string      `json:"scope"`
	Sid         string      `json:"sid"`
	Username    string      `json:"username"`
}

// CustomerClaims สำหรับ JWT ลูกค้า (ใช้ HS256)
type CustomerClaims struct {
	Name       string   `json:"name"`
	CustomerID uint     `json:"customer_id"`
	Email      string   `json:"email,omitempty"`
	Roles      []string `json:"roles"`
	jwt.RegisteredClaims
}

type RealmAccess struct {
	Roles []string `json:"roles"`
}
