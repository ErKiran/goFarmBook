package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/ErKiran/node/domain"
	"github.com/dgrijalva/jwt-go"
)

type authResponse struct {
	User  *domain.User     `json:"user"`
	Token *domain.JWTToken `json:"token"`
}

func (s *Server) registerUser() http.HandlerFunc {
	payload := domain.RegisterPayload{}
	return validatePayload(func(w http.ResponseWriter, r *http.Request) {
		user, err := s.domain.Register(payload)

		if err != nil {
			badRequestResponse(w, err)
			return
		}

		token, err := user.GenerateToken()

		if err != nil {
			badRequestResponse(w, err)
			return
		}

		jsonResponse(w, &authResponse{
			User:  user,
			Token: token,
		}, http.StatusCreated)
	}, &payload)
}

func (s *Server) withUser(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token, err := domain.ParseToken(r)

		fmt.Println("token", token)

		if err != nil {
			fmt.Println("In Error", err)
			unAuthorizedResponse(w)
			return
		}

		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			userID := int64(claims["id"].(float64))
			user, err := s.domain.GetUserById(userID)

			if err != nil {
				unAuthorizedResponse(w)
				return
			}

			ctx := context.WithValue(r.Context(), "currentUser", user)

			next.ServeHTTP(w, r.WithContext(ctx))

		} else {
			unAuthorizedResponse(w)
			return
		}
	})
}

func jsonResponse(w http.ResponseWriter, data interface{}, statusCode int) {
	w.Header().Set("Content-Type", "application/json")

	w.WriteHeader(statusCode)

	if data == nil {
		data = map[string]string{}
	}

	if err := json.NewEncoder(w).Encode(data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	return
}

func badRequestResponse(w http.ResponseWriter, err error) {
	response := map[string]string{"error": err.Error()}
	jsonResponse(w, response, http.StatusBadRequest)
}

func unAuthorizedResponse(w http.ResponseWriter) {
	response := map[string]string{"error": "UnAuthorized"}
	jsonResponse(w, response, http.StatusUnauthorized)
}

type PayloadValidation interface {
	IsValid() (bool, map[string]string)
}

func validatePayload(next http.HandlerFunc, payload PayloadValidation) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		err := json.NewDecoder(r.Body).Decode(&payload)
		if err != nil {
			badRequestResponse(w, err)
			return
		}

		defer r.Body.Close()

		if isValid, errs := payload.IsValid(); !isValid {
			jsonResponse(w, errs, http.StatusBadRequest)
			return
		}

		ctx := context.WithValue(r.Context(), "payload", payload)

		next.ServeHTTP(w, r.WithContext(ctx))

	}
}
