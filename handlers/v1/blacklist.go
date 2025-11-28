package v1

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// HandlePostBlacklist обрабатывает POST-запросы для добавления в чёрный список
func HandlePostBlacklist(db DatabaseInterface) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Only POST method is allowed", http.StatusMethodNotAllowed)
			return
		}

		w.Header().Set("Content-Type", "application/json")

		// Парсинг JSON из тела запроса
		var entry map[string]string
		err := json.NewDecoder(r.Body).Decode(&entry)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]string{
				"error": "Invalid JSON format",
			})
			return
		}

		socialNetwork := entry["social_network"]
		userID := entry["user_id"]
		nickname := entry["nickname"]

		// Валидация данных
		if !SocialNetworks[socialNetwork] {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(ResponsePayload{
				Status:  "error",
				Message: "Invalid social network",
			})
			return
		}

		if userID == "" {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(ResponsePayload{
				Status:  "error",
				Message: "User ID is required",
			})
			return
		}

		if nickname == "" {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(ResponsePayload{
				Status:  "error",
				Message: "Nickname is required",
			})
			return
		}

		// Добавление в чёрный список
		id, err := db.AddToBlacklist(socialNetwork, userID, nickname)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(ResponsePayload{
				Status:  "error",
				Message: "Failed to add user to blacklist: " + err.Error(),
			})
			return
		}

		response := ResponsePayload{
			Status:  "success",
			Message: fmt.Sprintf("User %s added to blacklist in %s", nickname, socialNetwork),
			Data: map[string]interface{}{
				"id":             id,
				"social_network": socialNetwork,
				"user_id":        userID,
				"nickname":       nickname,
			},
		}

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(response)
	}
}

// HandleGetBlacklist получает список всех забаненных пользователей
func HandleGetBlacklist(db DatabaseInterface) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "Only GET method is allowed", http.StatusMethodNotAllowed)
			return
		}

		w.Header().Set("Content-Type", "application/json")

		socialNetwork := r.URL.Query().Get("social_network")
		userID := r.URL.Query().Get("user_id")

		if socialNetwork == "" {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(ResponsePayload{
				Status:  "error",
				Message: "social_network parameter is required",
			})
			return
		}

		if userID == "" {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(ResponsePayload{
				Status:  "error",
				Message: "user_id parameter is required",
			})
			return
		}

		if !SocialNetworks[socialNetwork] {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(ResponsePayload{
				Status:  "error",
				Message: "Invalid social network",
			})
			return
		}

		bannedUsers, err := db.GetBlacklist(socialNetwork, userID)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(ResponsePayload{
				Status:  "error",
				Message: "Failed to retrieve blacklist",
			})
			return
		}

		if bannedUsers == nil {
			bannedUsers = []map[string]interface{}{}
		}

		response := ResponsePayload{
			Status:  "success",
			Message: "Blacklist retrieved successfully",
			Data: map[string]interface{}{
				"social_network": socialNetwork,
				"user_id":        userID,
				"banned_count":   len(bannedUsers),
				"banned_users":   bannedUsers,
			},
		}

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(response)
	}
}

// HandleCheckBlacklist проверяет наличие пользователя в чёрном списке
func HandleCheckBlacklist(db DatabaseInterface) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "Only GET method is allowed", http.StatusMethodNotAllowed)
			return
		}

		w.Header().Set("Content-Type", "application/json")

		socialNetwork := r.URL.Query().Get("social_network")
		userID := r.URL.Query().Get("user_id")
		nickname := r.URL.Query().Get("nickname")

		if socialNetwork == "" {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(ResponsePayload{
				Status:  "error",
				Message: "social_network parameter is required",
			})
			return
		}

		if userID == "" {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(ResponsePayload{
				Status:  "error",
				Message: "user_id parameter is required",
			})
			return
		}

		if nickname == "" {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(ResponsePayload{
				Status:  "error",
				Message: "nickname parameter is required",
			})
			return
		}

		if !SocialNetworks[socialNetwork] {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(ResponsePayload{
				Status:  "error",
				Message: "Invalid social network",
			})
			return
		}

		isBanned, err := db.CheckBlacklist(socialNetwork, userID, nickname)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(ResponsePayload{
				Status:  "error",
				Message: "Failed to check blacklist",
			})
			return
		}

		response := ResponsePayload{
			Status:  "success",
			Message: "Check completed",
			Data: map[string]interface{}{
				"social_network": socialNetwork,
				"user_id":        userID,
				"nickname":       nickname,
				"is_banned":      isBanned,
			},
		}

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(response)
	}
}

// HandleRoot обрабатывает GET-запросы к корню
func HandleRoot(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"message": "Welcome to Warden Solver Core",
		"version": "1.0.0",
	})
}
