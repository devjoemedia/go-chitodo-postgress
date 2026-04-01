package handlers

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/devjoemedia/go-ticketing-api/database"
	"github.com/devjoemedia/go-ticketing-api/models"
	api_response "github.com/devjoemedia/go-ticketing-api/types/response"
	"github.com/devjoemedia/go-ticketing-api/utils"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

var validate = validator.New()

// Register

// Register godoc
// @Summary      Register a new user
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        body body models.RegisterRequest true "Register payload"
// @Success      201 {object} api_response.RegisterResponse
// @Router       /api/v1/auth/register [post]
func Register(w http.ResponseWriter, r *http.Request) {
	var req models.RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.Error(w, http.StatusBadRequest, "Invalid JSON")
		return
	}
	if err := validate.Struct(&req); err != nil {
		utils.Error(w, http.StatusBadRequest, "Validation failed")
		return
	}

	// Check email uniqueness
	var existing models.User
	if err := database.DB.Where("email = ?", req.Email).First(&existing).Error; err == nil {
		utils.Error(w, http.StatusConflict, "Email already in use")
		return
	}

	// Hash password
	hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		utils.Error(w, http.StatusInternalServerError, "Failed to process password")
		return
	}

	user := models.User{
		Name:     req.Name,
		Email:    req.Email,
		Password: string(hash),
	}
	if error := database.DB.WithContext(r.Context()).Create(&user).Error; error != nil {
		utils.Error(w, http.StatusInternalServerError, "Failed to create user")
		return
	}

	tokens, err := issueTokenPair(r, &user)
	if err != nil {
		utils.Error(w, http.StatusInternalServerError, "Failed to issue tokens")
		return
	}

	utils.JSON(w, http.StatusCreated, api_response.RegisterResponse{
		Success: true,
		Status:  http.StatusCreated,
		Message: "Registration successful",
		User:    &user,
		Tokens:  tokens,
	})
}

// Login

// Login godoc
// @Summary      Login
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        body body models.LoginRequest true "Login payload"
// @Success      200 {object} api_response.LoginResponse
// @Router       /api/v1/auth/login [post]
func Login(w http.ResponseWriter, r *http.Request) {
	var req models.LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.Error(w, http.StatusBadRequest, "Invalid JSON")
		return
	}
	if err := validate.Struct(&req); err != nil {
		utils.Error(w, http.StatusBadRequest, "Validation failed")
		return
	}

	var user models.User
	if err := database.DB.Where("email = ?", req.Email).First(&user).Error; err != nil {
		// Deliberately vague — never reveal whether email exists
		utils.Error(w, http.StatusUnauthorized, "Invalid credentials")
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		utils.Error(w, http.StatusUnauthorized, "Invalid credentials")
		return
	}

	tokens, err := issueTokenPair(r, &user)
	if err != nil {
		utils.Error(w, http.StatusInternalServerError, "Failed to issue tokens")
		return
	}

	utils.JSON(w, http.StatusOK, api_response.LoginResponse{
		Success: true,
		Status:  http.StatusOK,
		Message: "Login successful",
		User:    &user,
		Tokens:  tokens,
	})
}

// Refresh

// RefreshToken godoc
// @Summary      Refresh access token
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        body body models.RefreshRequest true "Refresh token"
// @Success      200 {object} api_response.RefreshResponse
// @Router       /api/v1/auth/refresh [post]
func RefreshToken(w http.ResponseWriter, r *http.Request) {
	var req models.RefreshRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.Error(w, http.StatusBadRequest, "Invalid JSON")
		return
	}

	// Parse & validate JWT signature/expiry first
	claims, err := utils.ParseRefreshToken(req.RefreshToken)
	if err != nil {
		utils.Error(w, http.StatusUnauthorized, "Invalid or expired refresh token")
		return
	}

	// Lookup the DB record by the UUID embedded in the JWT
	var stored models.RefreshToken
	if lookup_error := database.DB.Where("token = ?", claims.TokenID).First(&stored).Error; lookup_error != nil {
		utils.Error(w, http.StatusUnauthorized, "Refresh token not found")
		return
	}

	// Guard: revoked or expired in DB
	if !stored.IsValid() {
		utils.Error(w, http.StatusUnauthorized, "Refresh token has been revoked or expired")
		return
	}

	// Rotation: revoke old, issue new pair
	stored.Revoked = true
	database.DB.Save(&stored)

	var user models.User
	if error := database.DB.First(&user, claims.UserID).Error; error != nil {
		utils.Error(w, http.StatusUnauthorized, "User not found")
		return
	}

	tokens, err := issueTokenPair(r, &user)
	if err != nil {
		utils.Error(w, http.StatusInternalServerError, "Failed to issue tokens")
		return
	}

	utils.JSON(w, http.StatusOK, api_response.RefreshResponse{
		Success: true,
		Status:  http.StatusOK,
		Message: "Token refreshed successfully",
		Tokens:  tokens,
	})
}

// Logout

// Logout godoc
// @Summary      Logout (revoke refresh token)
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        body body models.RefreshRequest true "Refresh token to revoke"
// @Success      200 {object} api_response.LogoutResponse
// @Router       /api/v1/auth/logout [post]
func Logout(w http.ResponseWriter, r *http.Request) {
	var req models.RefreshRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.Error(w, http.StatusBadRequest, "Invalid JSON")
		return
	}

	claims, err := utils.ParseRefreshToken(req.RefreshToken)
	if err != nil {
		// Token is already invalid — treat as successful logout
		utils.JSON(w, http.StatusOK, api_response.LogoutResponse{
			Success: true, Status: http.StatusOK, Message: "Logged out",
		})
		return
	}

	database.DB.Model(&models.RefreshToken{}).
		Where("token = ? AND user_id = ?", claims.TokenID, claims.UserID).
		Update("revoked", true)

	utils.JSON(w, http.StatusOK, api_response.LogoutResponse{
		Success: true,
		Status:  http.StatusOK,
		Message: "Logged out successfully",
	})
}

// Private helper

// issueTokenPair creates a DB-persisted refresh token and returns
// both a signed access JWT and a signed refresh JWT.
func issueTokenPair(r *http.Request, user *models.User) (api_response.AuthTokens, error) {
	tokenID := uuid.NewString()

	expiry := utils.RefreshExpiry()
	rt := models.RefreshToken{
		UserID:    user.ID,
		Token:     tokenID,
		ExpiresAt: time.Now().Add(expiry),
		UserAgent: r.UserAgent(),
		IPAddress: r.RemoteAddr,
	}
	if err := database.DB.Create(&rt).Error; err != nil {
		return api_response.AuthTokens{}, err
	}

	accessToken, err := utils.GenerateAccessToken(user.ID, user.Email)
	if err != nil {
		return api_response.AuthTokens{}, err
	}

	refreshToken, err := utils.GenerateRefreshToken(user.ID, tokenID)
	if err != nil {
		return api_response.AuthTokens{}, err
	}

	return api_response.AuthTokens{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		TokenType:    "Bearer",
		ExpiresIn:    15 * 60, // 15 minutes in seconds
	}, nil
}
