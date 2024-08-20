package http

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"time"

	"go-students-api/internal/database"
	"log"

	"github.com/dgrijalva/jwt-go"
)

type AuthClaims struct {
	UserID string `json:"user_id"`
	jwt.StandardClaims
}

func Login(w http.ResponseWriter, r *http.Request, jwtSecret []byte) {
	var loginRequest struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	if err := json.NewDecoder(r.Body).Decode(&loginRequest); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	// Fetch user from database
	db := database.GetDB()
	var user struct {
		ID       string
		Password string
	}
	err := db.QueryRow("SELECT id, password FROM users WHERE username = ?", loginRequest.Username).Scan(&user.ID, &user.Password)
	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		} else {
			log.Printf("Database query error: %v", err)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
		}
		return
	}

	// Validate password
	if user.Password != loginRequest.Password {
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}

	// Generate JWT token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &AuthClaims{
		UserID: user.ID,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour * 24).Unix(),
		},
	})

	tokenString, err := token.SignedString(jwtSecret)
	if err != nil {
		log.Printf("Error generating token: %v", err)
		http.Error(w, "Failed to generate token", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"token": tokenString})
}
