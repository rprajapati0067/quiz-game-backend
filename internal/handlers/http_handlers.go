package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/rprajapati0067/quiz-app-tools/logger"
	"github.com/rprajapati0067/quiz-game-backend/internal/service"
)

type HTTPHandlers struct {
	authService     service.AuthService
	userService     service.UserService
	questionService service.QuestionService
}

func NewHTTPHandlers(authService service.AuthService, userService service.UserService, questionService service.QuestionService) *HTTPHandlers {
	return &HTTPHandlers{
		authService:     authService,
		userService:     userService,
		questionService: questionService,
	}
}

func (h *HTTPHandlers) Health(w http.ResponseWriter, r *http.Request) {
	logger.Info("Health check request received")
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"status":  "ok",
		"service": "quiz-backend",
	})
}

func (h *HTTPHandlers) Signup(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req struct {
		Name  string `json:"name"`
		Phone string `json:"phone"`
		Email string `json:"email"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	user, err := h.authService.Signup(r.Context(), req.Name, req.Phone, req.Email)
	if err != nil {
		logger.Error(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"user_id": user.ID})
}

func (h *HTTPHandlers) Login(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req struct {
		Phone string `json:"phone"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	user, err := h.authService.Login(r.Context(), req.Phone)
	if err != nil {
		logger.Error(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if user == nil {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"token": "dummy-token"}) // TODO: generate real token
}

func (h *HTTPHandlers) Me(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// TODO: Extract user ID from JWT token
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"user_id": "",
		"name":    "",
		"phone":   "",
		"email":   "",
		"points":  0,
	})
}

func (h *HTTPHandlers) ListQuestions(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	slotStr := r.URL.Query().Get("slot")
	slot, err := strconv.ParseInt(slotStr, 10, 32)
	if err != nil || slot <= 0 {
		http.Error(w, "Invalid slot parameter", http.StatusBadRequest)
		return
	}

	questions, err := h.questionService.ListBySlot(r.Context(), int32(slot))
	if err != nil {
		logger.Error(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(questions)
}

func (h *HTTPHandlers) CreateQuestion(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req struct {
		Text         string   `json:"text"`
		Options      []string `json:"options"`
		CorrectIndex int32    `json:"correct_index"`
		Slot         int32    `json:"slot"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	question, err := h.questionService.Create(r.Context(), req.Text, req.Options, req.CorrectIndex, req.Slot, "admin")
	if err != nil {
		logger.Error(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(question)
}

func (h *HTTPHandlers) SubmitAnswer(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// TODO: Implement answer submission
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"correct":        false,
		"updated_points": 0,
	})
}

func (h *HTTPHandlers) SetupRoutes(mux *http.ServeMux) {
	// Health endpoints
	mux.HandleFunc("/health", h.Health)
	mux.HandleFunc("/healthz", h.Health)

	// Auth endpoints
	mux.HandleFunc("/api/v1/auth/signup", h.Signup)
	mux.HandleFunc("/api/v1/auth/login", h.Login)

	// User endpoints
	mux.HandleFunc("/api/v1/user/me", h.Me)

	// Question endpoints
	mux.HandleFunc("/api/v1/questions", h.ListQuestions)
	mux.HandleFunc("/api/v1/questions/create", h.CreateQuestion)
	mux.HandleFunc("/api/v1/questions/submit", h.SubmitAnswer)
}

