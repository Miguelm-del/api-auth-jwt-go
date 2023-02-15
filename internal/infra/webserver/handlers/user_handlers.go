package handlers

import (
	"encoding/json"
	"github.com/Miguelm-del/api-auth-jwt-go/internal/dto"
	"github.com/Miguelm-del/api-auth-jwt-go/internal/entity"
	"github.com/Miguelm-del/api-auth-jwt-go/internal/infra/database"
	"github.com/go-chi/jwtauth"
	"net/http"
	"time"
)

type UserHandler struct {
	UserDB database.UserInterface
}

type Error struct {
	Message string `json:"message"`
}

func NewUserHandler(userDB database.UserInterface) *UserHandler {
	return &UserHandler{
		UserDB: userDB,
	}
}

// GetJWT user godoc
// @Summary Get a user JWT
// @Description Get a user JWT
// @Tags Users
// @Accept json
// @Produce json
// @Param request body dto.GetJWTInput true "User request"
// @Success 200 {object} dto.GetJWTOutput
// Failure 404 {object} Error
// @Failure 500 {object} Error
// @Router /users/generate_token [post]
func (h *UserHandler) GetJWT(w http.ResponseWriter, r *http.Request) {
	jwt := r.Context().Value("jwt").(*jwtauth.JWTAuth)
	JwtExpiresIn := r.Context().Value("JwtExpiresIn").(int)
	var user dto.GetJWTInput
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		errorMsg := Error{Message: err.Error()}
		json.NewEncoder(w).Encode(errorMsg)
		return
	}
	u, err := h.UserDB.FindByEmail(user.Email)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		errorMsg := Error{Message: err.Error()}
		json.NewEncoder(w).Encode(errorMsg)
		return
	}
	if !u.ValidatePassword(user.Password) {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	_, tokenString, _ := jwt.Encode(map[string]interface{}{
		"sub": u.ID.String(),
		"exp": time.Now().Add(time.Second * time.Duration(JwtExpiresIn)).Unix(),
	})

	accessToken := dto.GetJWTOutput{AccessToken: tokenString}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(accessToken)
}

// Create user godoc
// @Summary Create User
// @Description Create user
// @Tags Users
// @Accept json
// @Produce json
// @Param request body dto.CreateUserInput true "User request"
// @Success 201
// @Failure 500 {object} Error
// @Router /users [post]
func (h *UserHandler) Create(w http.ResponseWriter, r *http.Request) {
	var user dto.CreateUserInput
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	u, err := entity.NewUser(user.Name, user.Email, user.Password)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		errorMsg := Error{Message: err.Error()}
		json.NewEncoder(w).Encode(errorMsg)
		return
	}
	err = h.UserDB.Create(u)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		errorMsg := Error{Message: err.Error()}
		json.NewEncoder(w).Encode(errorMsg)
		return
	}
	w.WriteHeader(http.StatusCreated)
}
