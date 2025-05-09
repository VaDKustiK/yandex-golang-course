package orchestrator

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/VaDKustiK/yandex-golang-course/calculator_service/pkg/common"
	"golang.org/x/crypto/bcrypt"

	"github.com/golang-jwt/jwt/v4"
)

var jwtKey = []byte("super_secret_key")

type Claims struct {
	UserID uint `json:"user_id"`
	jwt.RegisteredClaims
}

type credentials struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

func RegisterHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	var creds credentials
	if err := json.NewDecoder(r.Body).Decode(&creds); err != nil {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(creds.Password), bcrypt.DefaultCost)
	if err != nil {
		http.Error(w, "could not hash password", http.StatusInternalServerError)
		return
	}

	user := common.User{Login: creds.Login, Password: string(hash)}
	if err := GetDB().Create(&user).Error; err != nil {
		http.Error(w, "cannot create user", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	var creds credentials
	if err := json.NewDecoder(r.Body).Decode(&creds); err != nil {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}

	var user common.User
	if err := GetDB().Where("login = ?", creds.Login).First(&user).Error; err != nil {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}
	if bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(creds.Password)) != nil {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	claims := &Claims{
		UserID: user.ID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		http.Error(w, "could not sign token", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(map[string]string{"token": tokenString})
}
