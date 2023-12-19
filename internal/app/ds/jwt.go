package ds

import (
	"R_I_P_labs/internal/app/role"
	"github.com/golang-jwt/jwt"
)

type JWTClaims struct {
	jwt.StandardClaims
	UserUUID string `json:"user_uuid"`
	Role     role.Role
}