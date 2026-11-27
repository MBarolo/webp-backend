package middleware

import (
	"context"
	"crypto/subtle"
	"encoding/base64"
	"errors"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"github.com/mbarolo/webp_backend/models"
	"github.com/mbarolo/webp_backend/utils"
)

type Claims struct {
	Sub       string `json:"sub"`
	Exp       int64  `json:"exp"`
	Email     string `json:"email"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
}

// Valid implementa la interfaz jwt.Claims
func (c Claims) Valid() error {
	if c.Exp < time.Now().Unix() {
		return errors.New("token expirado")
	}
	return nil
}

const EXPIRATION_IN_DAYS = 30

var jwtKey []byte

func init() {
	// Asignamos la key del .env, si no existe se asigna una default
	key := os.Getenv("JWT_KEY")
	if key == "" {
		key = "default_key"
	}

	jwtKey = []byte(key)
}

// GenerateToken: Genera un token json segun las claims del struct
func GenerateToken(user models.User) (string, time.Time, error) {
	expirationTime := time.Now().Add(time.Hour * 24 * EXPIRATION_IN_DAYS)
	claims := &Claims{
		Sub:       fmt.Sprintf("%d", user.Id),
		Exp:       expirationTime.Unix(),
		Email:     user.Email,
		FirstName: user.FirstName,
		LastName:  user.LastName,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		return "", time.Time{}, err
	}

	return tokenString, expirationTime, nil
}

// validateToken: Valida un token json
func validateToken(tokenString string) (*Claims, error) {
	claims := &Claims{}

	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("método de firma inválido")
		}
		return jwtKey, nil
	})

	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, errors.New("token inválido")
	}

	return claims, nil
}

// Middleware que valida el token enviado, si no es encontrado retorna 401
func AuthMiddleware(c *gin.Context) {
	header := c.GetHeader("Authorization")
	token := ""

	if header != "" {
		parts := strings.Split(header, " ")
		if len(parts) == 2 && parts[0] == "Bearer" {
			token = parts[1]
		} else {
			utils.JsonResponse(c, http.StatusUnauthorized, "Estructura de header de autorización inesperada", nil)
			c.Abort()
			return
		}
	} else {
		utils.JsonResponse(c, http.StatusUnauthorized, "Token no proporcionado", nil)
		c.Abort()
		return
	}

	claims, err := validateToken(token)
	if err != nil {
		utils.JsonResponse(c, http.StatusUnauthorized, "Error al validar el token: "+err.Error(), nil)
		c.Abort()
		return
	}

	ctx := context.WithValue(c.Request.Context(), "claims", claims)
	c.Request = c.Request.WithContext(ctx)
	c.Next()
}

// validateAdminCredentials compara las credenciales del encabezado de autorización con las del .env
func validateAdminCredentials(creds string) (bool, error) {
	adminCredentials, ok := os.LookupEnv("ADMIN_CREDENTIALS")
	if !ok {
		return false, errors.New("credenciales de administrador no presentes en el entorno")
	}

	if _, err := base64.StdEncoding.DecodeString(creds); err != nil {
		return false, errors.New("credenciales no son un base64 válido")
	}

	return subtle.ConstantTimeCompare([]byte(adminCredentials), []byte(creds)) == 1, nil
}

// Middleware para rutas admin
func AdminMiddleware(c *gin.Context) {
	header := c.GetHeader("Authorization")
	creds := ""
	if header != "" {
		parts := strings.Split(header, " ")
		if len(parts) == 2 && parts[0] == "Basic" {
			creds = parts[1]
		} else {
			utils.JsonResponse(c, http.StatusUnauthorized, "Estructura de header de autorización inesperada", nil)
			c.Abort()
			return
		}
	}

	if valid, err := validateAdminCredentials(creds); err != nil {
		utils.JsonResponse(c, http.StatusUnauthorized, "Error al validar las credenciales"+err.Error(), nil)
		c.Abort()
		return
	} else if !valid {
		utils.JsonResponse(c, http.StatusUnauthorized, "Credenciales inválidas", nil)
		c.Abort()
		return
	}

	ctx := context.WithValue(c.Request.Context(), "claims", "")
	c.Request = c.Request.WithContext(ctx)
	c.Next()
}
