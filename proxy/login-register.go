package main

import (
	"encoding/json"
	"github.com/dgrijalva/jwt-go"
	"golang.org/x/crypto/bcrypt"
	"net/http"
)

type User struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

var users = make(map[string]string)

// register обрабатывает процесс регистрации пользователя
//
// @Summary Регистрация пользователя
// @Description Позволяет пользователям зарегистрироваться и получить JWT-токен.
// @Tags auth
// @Accept json
// @Produce json
// @Param user body User true "Информация о пользователе"
// @Success 200 {object} map[string]string "Успешная регистрация, возвращает сообщение"
// @Failure 400 {string} string "Неверный формат запроса"
// @Failure 500 {string} string "Внутренняя ошибка сервера"
// @Router /api/register [post]
func register(w http.ResponseWriter, r *http.Request) {
	var user User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, "Invalid request format", http.StatusBadRequest)
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	users[user.Username] = string(hashedPassword)

	w.WriteHeader(http.StatusOK)
}

// login обрабатывает процесс входа пользователя
//
// @Summary Вход пользователя
// @Description Позволяет пользователям войти в систему и получить JWT-токен.
// @Tags auth
// @Accept json
// @Produce json
// @Param input body User true "Информация о пользователе для входа"
// @Success 200 {object} map[string]string "Успешный вход, возвращает JWT-токен"
// @Failure 400 {string} string "Неверный формат запроса"
// @Failure 401 {string} string "Пользователь не найден или неверные учетные данные"
// @Failure 500 {string} string "Внутренняя ошибка сервера"
// @Router /api/login [post]
func login(w http.ResponseWriter, r *http.Request) {
	var user User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, "Invalid request format", http.StatusBadRequest)
		return
	}

	storedPassword, exists := users[user.Username]
	if !exists {
		http.Error(w, "User not found", http.StatusUnauthorized)
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(storedPassword), []byte(user.Password))
	if err != nil {
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}

	claims := jwt.MapClaims{"username": user.Username}
	_, tokenString, err := tokenAuth.Encode(claims)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"token": tokenString})
}
