package handler

import (
	"YandexLearnMiddle/internal/store"
	"YandexLearnMiddle/postgresql"
	"database/sql"
	"encoding/json"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
	"net/http"
)

var Logger *zap.Logger

func SignUp(w http.ResponseWriter, r *http.Request) {
	var user store.User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		Logger.Error("Invalid request body", zap.Error(err))
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	var exists bool
	err := postgresql.DB.QueryRow("SELECT EXISTS(SELECT 1 FROM users WHERE username=$1)", user.Username).Scan(&exists)
	if err != nil {
		Logger.Error("Database error during user existence check", zap.Error(err))
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}

	if exists {
		Logger.Warn("User already exists", zap.String("username", user.Username))
		http.Error(w, "User already exists", http.StatusConflict)
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		Logger.Error("Error while hashing password", zap.Error(err))
		http.Error(w, "Error while hashing password", http.StatusInternalServerError)
		return
	}

	_, err = postgresql.DB.Exec("INSERT INTO users (username, password_hash) VALUES ($1, $2)", user.Username, string(hashedPassword))
	if err != nil {
		Logger.Error("Database error during user registration", zap.Error(err))
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}

	Logger.Info("User registered successfully", zap.String("username", user.Username))
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte("User registered successfully"))
}

func SignIn(w http.ResponseWriter, r *http.Request) {
	var user store.User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		Logger.Error("Invalid request body", zap.Error(err))
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	var storedPasswordHash string
	err := postgresql.DB.QueryRow("SELECT password_hash FROM users WHERE username=$1", user.Username).Scan(&storedPasswordHash)
	if err != nil {
		if err == sql.ErrNoRows {
			Logger.Warn("Invalid login attempt: user not found", zap.String("username", user.Username))
			http.Error(w, "Invalid username or password", http.StatusUnauthorized)
			return
		}
		Logger.Error("Database error during login", zap.Error(err))
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(storedPasswordHash), []byte(user.Password)); err != nil {
		Logger.Warn("Invalid login attempt: wrong password", zap.String("username", user.Username))
		http.Error(w, "Invalid username or password", http.StatusUnauthorized)
		return
	}

	token, err := CreateJWT(user.Username)
	if err != nil {
		Logger.Error("Failed to create JWT token", zap.Error(err))
		http.Error(w, "Could not create token", http.StatusInternalServerError)
		return
	}

	Logger.Info("User logged in successfully", zap.String("username", user.Username))
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"token": token,
	})
}
