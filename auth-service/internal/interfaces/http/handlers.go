package http

import (
	//"context"
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	authapp "github.com/oziev02/taskflow-microservices/auth-service/internal/application/auth"
	jwtmgr "github.com/oziev02/taskflow-microservices/auth-service/internal/infrastructure/jwt"
)

type Handler struct {
	svc *authapp.Service
	jwt *jwtmgr.Manager
}

func RegisterRoutes(r chi.Router, svc *authapp.Service, jwt *jwtmgr.Manager) {
	h := &Handler{svc: svc, jwt: jwt}

	r.Route("/auth", func(r chi.Router) {
		r.Post("/register", h.Register)
		r.Post("/login", h.Login)
		r.Post("/refresh", h.Refresh)

		// Пример защищенного эндпойнта
		r.Group(func(priv chi.Router) {
			priv.Use(AuthMiddleware(jwt))
			priv.Get("/me", h.Me)
		})
	})
}

func (h *Handler) Register(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil || req.Email == "" || req.Password == "" {
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}
	id, err := h.svc.Register(r.Context(), req.Email, req.Password)
	if err != nil {
		http.Error(w, err.Error(), http.StatusConflict)
		return
	}
	_ = json.NewEncoder(w).Encode(map[string]any{"user_id": id})
}

func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}
	at, rt, err := h.svc.Login(r.Context(), req.Email, req.Password)
	if err != nil {
		http.Error(w, "invalid credentials", http.StatusUnauthorized)
		return
	}
	_ = json.NewEncoder(w).Encode(map[string]string{
		"access_token":  at,
		"refresh_token": rt,
	})
}

func (h *Handler) Refresh(w http.ResponseWriter, r *http.Request) {
	var req struct {
		RefreshToken string `json:"refresh_token"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil || req.RefreshToken == "" {
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}
	at, rt, err := h.svc.Refresh(r.Context(), req.RefreshToken)
	if err != nil {
		http.Error(w, "refresh failed", http.StatusUnauthorized)
		return
	}
	_ = json.NewEncoder(w).Encode(map[string]string{
		"access_token":  at,
		"refresh_token": rt,
	})
}

func (h *Handler) Me(w http.ResponseWriter, r *http.Request) {
	uid := r.Context().Value(userIDKey)
	_ = json.NewEncoder(w).Encode(map[string]any{"user_id": uid})
}
