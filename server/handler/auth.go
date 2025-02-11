package handler

import (
	"log"
	"server/middleware"
	"server/model"

	"github.com/go-playground/validator"
	"github.com/jmoiron/sqlx"
	"github.com/labstack/echo/v4"
	"golang.org/x/crypto/bcrypt"

	"github.com/dgrijalva/jwt-go"
)

const (
	DEFAULT_ICON_URL = "http://localhost:5173/images/default-icon.webp"
	DEFAULT_HEADER_URL = "http://localhost:5173/images/default-header.webp"
	DEFAULT_PROFILE = ""
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
	ID        int    `json:"id"`
	Name      string `json:"name"`
	DisplayID string `json:"display_id"`
	IconURL   string `json:"icon_url"`
}

type TokenResponse struct {
	Token string            `json:"token"`
	User  TokenUserResponse `json:"user"`
}

func (h *AuthHandler) Signup(c echo.Context) error {
	req := new(SignupRequest)
	if err := c.Bind(req); err != nil {
		log.Printf("%v", err)
		return c.JSON(400, map[string]string{"message": "Bad Request"})
	}

	if err := h.validator.Struct(req); err != nil {
		log.Printf("%v", err)
		return c.JSON(400, map[string]string{"message": "Bad Request"})
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		log.Printf("%v", err)
		return c.JSON(500, map[string]string{"message": "Internal Server Error"})
	}

	res, err := h.db.Exec("INSERT INTO users (email, password_hash) VALUES (?, ?)", req.Email, string(hash))
	if err != nil {
		log.Printf("%v", err)
		return c.JSON(500, map[string]string{"message": "Internal Server Error"})
	}

	id, err := res.LastInsertId()
	if err != nil {
		log.Printf("%v", err)
		return c.JSON(500, map[string]string{"message": "Internal Server Error"})
	}
	_, err = h.db.Exec("INSERT INTO user_profiles (user_id, name, display_id, header_url, icon_url, profile) VALUES (?, ?, ?, ?, ?, ?)", id, req.Name, req.DisplayID, DEFAULT_HEADER_URL, DEFAULT_ICON_URL, DEFAULT_PROFILE)

	token, err := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{"user_id": id}).SignedString([]byte(h.secret))
	if err != nil {
		log.Printf("%v", err)
		return c.JSON(500, map[string]string{"message": "Internal Server Error"})
	}

	var user model.User
	err = h.db.Get(&user, "SELECT * FROM users WHERE id = ?", id)
	if err != nil {
		log.Printf("%v", err)
		return c.JSON(500, map[string]string{"message": "Internal Server Error"})
	}

	userProfile, err := h.getUserProfileByID(int(id))
	if err != nil {
		log.Printf("%v", err)
		return c.JSON(500, map[string]string{"message": "Internal Server Error"})
	}

	log.Printf("%v", userProfile)

	return c.JSON(200, TokenResponse{Token: token, User: TokenUserResponse{ID: user.ID, Name: userProfile.Name, DisplayID: userProfile.DisplayID, IconURL: userProfile.IconURL}})
}

type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8"`
}

func (h *AuthHandler) Login(c echo.Context) error {
	req := new(LoginRequest)
	if err := c.Bind(req); err != nil {
		log.Printf("%v", err)
		return c.JSON(400, map[string]string{"message": "Bad Request"})
	}

	if err := h.validator.Struct(req); err != nil {
		log.Printf("%v", err)
		return c.JSON(400, map[string]string{"message": "Bad Request"})
	}

	var user model.User
	err := h.db.Get(&user, "SELECT * FROM users WHERE email = ?", req.Email)
	if err != nil {
		log.Printf("%v", err)
		return c.JSON(401, map[string]string{"message": "Unauthorized"})
	}

	userProfile, err := h.getUserProfileByID(user.ID)
	if err != nil {
		log.Printf("%v", err)
		return c.JSON(500, map[string]string{"message": "Internal Server Error"})
	}

	log.Printf("%v", userProfile)

	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password))
	if err != nil {
		log.Printf("%v", err)
		return c.JSON(401, map[string]string{"message": "Unauthorized"})
	}

	token, err := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{"user_id": user.ID}).SignedString([]byte(h.secret))
	if err != nil {
		log.Printf("%v", err)
		return c.JSON(500, map[string]string{"message": "Internal Server Error"})
	}

	return c.JSON(200, TokenResponse{Token: token, User: TokenUserResponse{ID: user.ID, Name: userProfile.Name, DisplayID: userProfile.DisplayID, IconURL: userProfile.IconURL}})
}

type MeResponse struct {
	ID        int    `json:"id"`
	Name      string `json:"name"`
	DisplayID string `json:"display_id"`
	IconURL   string `json:"icon_url"`
}

func (h *AuthHandler) Me(c echo.Context) error {
	userID := c.Get("user_id").(int)

	var user model.User
	err := h.db.Get(&user, "SELECT * FROM users WHERE id = ?", userID)
	if err != nil {
		return c.JSON(500, map[string]string{"message": "Internal Server Error"})
	}

	userProfile, err := h.getUserProfileByID(userID)
	if err != nil {
		return c.JSON(500, map[string]string{"message": "Internal Server Error"})
	}

	return c.JSON(200, MeResponse{ID: user.ID, Name: userProfile.Name, DisplayID: userProfile.DisplayID, IconURL: userProfile.IconURL})
}

type UserProfile struct {
	Name      string `db:"name"`
	DisplayID string `db:"display_id"`
	IconURL   string `db:"icon_url"`
}

func (h *AuthHandler) getUserProfileByID(userID int) (UserProfile, error) {
	var userProfile UserProfile
	// user_profilesテーブルからname, display_id, icon_urlを取得
	err := h.db.Get(&userProfile, "SELECT name, display_id, icon_url FROM user_profiles WHERE user_id = ?", userID)
	if err != nil {
		return userProfile, err
	}

	return userProfile, nil
}
