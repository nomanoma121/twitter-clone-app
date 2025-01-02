package middleware

import (
	"github.com/jmoiron/sqlx"
	"github.com/labstack/echo/v4"

	"github.com/dgrijalva/jwt-go"
)

type AuthMiddleware struct {
	db     *sqlx.DB
	secret string
}

func NewAuthMiddleware(db *sqlx.DB, secret string) *AuthMiddleware {
	return &AuthMiddleware{db: db, secret: secret}
}

func (m *AuthMiddleware) Middleware() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			token := c.Request().Header.Get("Authorization")
			if token == "" {
				return c.JSON(401, map[string]string{"message": "トークンが必要です"})
			}

			token = token[len("Bearer "):]

			t, err := jwt.Parse(token, func(t *jwt.Token) (interface{}, error) {
				if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
					return nil, jwt.ErrSignatureInvalid
				}
				return []byte(m.secret), nil
			})
			if err != nil {
				return c.JSON(401, map[string]string{"message": "トークンが無効です"})
			}

			if !t.Valid {
				return c.JSON(401, map[string]string{"message": "トークンが無効です"})
			}

			claims, ok := t.Claims.(jwt.MapClaims)
			if !ok {
				return c.JSON(401, map[string]string{"message": "トークンが無効です"})
			}

			userID, ok := claims["user_id"].(float64)
			if !ok {
				return c.JSON(401, map[string]string{"message": "トークンが無効です"})
			}

			c.Set("user_id", int(userID))

			return next(c)
		}
	}
}
