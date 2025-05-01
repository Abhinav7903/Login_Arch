package server

import (
	"LoginArch/factory"
	"LoginArch/pkg/users"
	"encoding/json"
	"log/slog"
	"net/http"
	"time"
)

func (s *Server) HandleCreateUser() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var user factory.User
		if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
			s.respond(
				w,
				ResponseMsg{
					Message: "Invalid request payload",
				},
				http.StatusBadRequest,
				nil,
			)
			return
		}

		//convert password to hash
		var err error
		user.HashedPassword, err = users.HashPassword(user.Password)
		if err != nil {
			s.respond(
				w,
				ResponseMsg{
					Message: "Error hashing password",
				},
				http.StatusInternalServerError,
				nil,
			)
			return
		}

		if err := s.user.CreateUser(user); err != nil {
			s.respond(
				w,
				ResponseMsg{
					Message: "Error creating user",
				},
				http.StatusInternalServerError,
				nil,
			)
			return
		}

		s.respond(
			w,
			ResponseMsg{
				Message: "User created successfully",
			},
			http.StatusCreated,
			nil,
		)
		slog.Debug("User created successfully", "user", user.Email)
	}
}

func (s *Server) HandleLogin() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var loginUser factory.User
		if err := json.NewDecoder(r.Body).Decode(&loginUser); err != nil {
			http.Error(w, "Invalid request payload", http.StatusBadRequest)
			return
		}

		user, err := s.user.Login(loginUser)
		if err != nil {
			http.Error(w, "Authentication failed", http.StatusUnauthorized)
			slog.Warn("Login failed", "email", loginUser.Email, "error", err)
			return
		}

		// Check if user has an active session
		sessionExists, err := s.redis.CheckSession(user.Email)
		if err != nil {
			http.Error(w, "In---ternal error", http.StatusInternalServerError)
			return
		}
		if sessionExists {
			http.Error(w, "User already logged in elsewhere", http.StatusConflict)
			return
		}

		// Create a session token
		token, err := s.redis.GenerateToken(user.Email)
		if err != nil {
			http.Error(w, "Internal error", http.StatusInternalServerError)
			return
		}

		// Store session in Redis
		if err := s.redis.StoreSession(user.Email, token); err != nil {
			http.Error(w, "Internal error", http.StatusInternalServerError)
			return
		}

		s.respond(
			w,
			ResponseMsg{
				Message: "Login successful",
				Data: map[string]string{
					"email": user.Email,
					"token": token,
				},
			},
			http.StatusOK,
			nil,
		)
	}
}

func (s *Server) HandleGetUser() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		email := r.URL.Query().Get("email")
		if email == "" {
			s.respond(
				w,
				ResponseMsg{
					Message: "Email parameter is required",
				},
				http.StatusBadRequest,
				nil,
			)
			return
		}

		user, err := s.user.GetUser(email)
		if err != nil {
			s.respond(
				w,
				ResponseMsg{
					Message: "Error retrieving user",
				},
				http.StatusInternalServerError,
				nil,
			)
			return
		}

		loc, err := time.LoadLocation("Asia/Kolkata")
		if err != nil {
			s.respond(
				w,
				ResponseMsg{
					Message: "Error loading location",
				},
				http.StatusInternalServerError,
				nil,
			)
			return
		}

		// Parse the created time (string) to time.Time
		createdTime, err := time.Parse(time.RFC3339, user.Created)
		if err != nil {
			s.respond(
				w,
				ResponseMsg{
					Message: "Error parsing created time",
				},
				http.StatusInternalServerError,
				nil,
			)
			return
		}

		// Format the time to desired format
		formattedCreated := createdTime.In(loc).Format("2006-01-02 15:04:05")

		// Update user.Created to formatted string
		user.Created = formattedCreated

		s.respond(
			w,
			ResponseMsg{
				Message: "success",
				Data:    user,
			},
			http.StatusOK,
			nil,
		)
	}
}

// LogoutHandler handles user logout
func (s *Server) HandleLogout() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		email := r.URL.Query().Get("email")
		if email == "" {
			s.respond(
				w,
				ResponseMsg{
					Message: "Email parameter is required",
				},
				http.StatusBadRequest,
				nil,
			)
			return
		}

		if err := s.redis.DeleteSession(email); err != nil {
			s.respond(
				w,
				ResponseMsg{
					Message: "Error deleting session",
				},
				http.StatusInternalServerError,
				nil,
			)
			return
		}

		s.respond(
			w,
			ResponseMsg{
				Message: "Logout successful",
			},
			http.StatusOK,
			nil,
		)
	}
}
