package handler

import (
	"server/middleware"
	"server/model"

	"github.com/go-playground/validator"
	"github.com/jmoiron/sqlx"
	"github.com/labstack/echo/v4"
	"golang.org/x/crypto/bcrypt"

	"github.com/dgrijalva/jwt-go"
)

type AuthConfig struct {
	Hash string
}

type AuthHandler struct {
	db        *sqlx.DB
	validator *validator.Validate
	secret    string
}

func NewAuthHandler(db *sqlx.DB, secret string) *AuthHandler {
	return &AuthHandler{db: db, validator: validator.New(), secret: secret}
}

func (h *AuthHandler) Register(g *echo.Group, authMiddleware *middleware.AuthMiddleware) {
	g.POST("/signup", h.Signup)
	g.POST("/login", h.Login)
	g.GET("/me", h.Me, authMiddleware.Middleware())
}

type SignupRequest struct {
	Name      string `json:"name" validate:"required"`
	Email     string `json:"email" validate:"required,email"`
	DisplayID string `json:"display_id" validate:"required"`
	Password  string `json:"password" validate:"required,min=8"`
}

type TokenUserResponse struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

type TokenResponse struct {
	Token string            `json:"token"`
	User  TokenUserResponse `json:"user"`
}

func (h *AuthHandler) Signup(c echo.Context) error {
	req := new(SignupRequest)
	if err := c.Bind(req); err != nil {
		return c.JSON(400, map[string]string{"message": "Bad Request"})
	}

	if err := h.validator.Struct(req); err != nil {
		return c.JSON(400, map[string]string{"message": "Bad Request"})
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return c.JSON(500, map[string]string{"message": "Internal Server Error"})
	}

	res, err := h.db.Exec("INSERT INTO users (email, password_hash) VALUES (?, ?)", req.Email, string(hash))
	if err != nil {
		return c.JSON(500, map[string]string{"message": "Internal Server Error"})
	}

	id, err := res.LastInsertId()
	if err != nil {
		return c.JSON(500, map[string]string{"message": "Internal Server Error"})
	}
	_, err = h.db.Exec("INSERT INTO user_profiles (user_id, name, display_id) VALUES (?, ?, ?)", id, req.Name, req.DisplayID)

	token, err := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{"user_id": id}).SignedString([]byte(h.secret))
	if err != nil {
		return c.JSON(500, map[string]string{"message": "Internal Server Error"})
	}

	var user model.User
	err = h.db.Get(&user, "SELECT * FROM users WHERE id = ?", id)
	if err != nil {
		return c.JSON(500, map[string]string{"message": "Internal Server Error"})
	}
	var username string
	err = h.db.Get(&username, "SELECT name FROM user_profiles WHERE user_id = ?", user.ID)
	if err != nil {
		return c.JSON(500, map[string]string{"message": "Internal Server Error"})
	}

	return c.JSON(200, TokenResponse{Token: token, User: TokenUserResponse{ID: user.ID, Name: username, Email: user.Email}})
}

type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8"`
}

func (h *AuthHandler) Login(c echo.Context) error {
	req := new(LoginRequest)
	if err := c.Bind(req); err != nil {
		return c.JSON(400, map[string]string{"message": "Bad Request"})
	}

	if err := h.validator.Struct(req); err != nil {
		return c.JSON(400, map[string]string{"message": "Bad Request"})
	}

	var user model.User
	err := h.db.Get(&user, "SELECT * FROM users WHERE email = ?", req.Email)
	if err != nil {
		return c.JSON(401, map[string]string{"message": "Unauthorized"})
	}
	var username string
	err = h.db.Get(&username, "SELECT name FROM user_profiles WHERE user_id = ?", user.ID)
	if err != nil {
		return c.JSON(401, map[string]string{"message": "Unauthorized"})
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password))
	if err != nil {
		return c.JSON(401, map[string]string{"message": "Unauthorized"})
	}

	token, err := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{"user_id": user.ID}).SignedString([]byte(h.secret))
	if err != nil {
		return c.JSON(500, map[string]string{"message": "Internal Server Error"})
	}

	return c.JSON(200, TokenResponse{Token: token, User: TokenUserResponse{ID: user.ID, Name: username, Email: user.Email}})
}

type MeResponse struct {
	ID        int    `json:"id"`
	Name      string `json:"name"`
	DisplayID string `json:"display_id"`
	Email     string `json:"email"`
}

func (h *AuthHandler) Me(c echo.Context) error {
	userID := c.Get("user_id").(int)

	var user model.User
	err := h.db.Get(&user, "SELECT * FROM users WHERE id = ?", userID)
	if err != nil {
		return c.JSON(500, map[string]string{"message": "Internal Server Error"})
	}
	var userProfile model.UserProfile
	err = h.db.Get(&userProfile, "SELECT name, display_id FROM user_profiles WHERE user_id = ?", userID)
	if err != nil {
		return c.JSON(500, map[string]string{"message": "Internal Server Error"})
	}

	return c.JSON(200, MeResponse{ID: user.ID, Name: userProfile.Name, Email: user.Email, DisplayID: userProfile.DisplayID})
}
