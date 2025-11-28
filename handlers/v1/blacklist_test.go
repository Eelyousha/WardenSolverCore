package v1

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

// MockDatabase реализует DatabaseInterface для тестирования
type MockDatabase struct {
	addToBlacklistFunc func(socialNetwork, userID, nickname string) (int, error)
	getBlacklistFunc   func(socialNetwork, userID string) ([]map[string]interface{}, error)
	checkBlacklistFunc func(socialNetwork, userID, nickname string) (bool, error)
}

func (m *MockDatabase) AddToBlacklist(socialNetwork, userID, nickname string) (int, error) {
	if m.addToBlacklistFunc != nil {
		return m.addToBlacklistFunc(socialNetwork, userID, nickname)
	}
	return 0, nil
}

func (m *MockDatabase) GetBlacklist(socialNetwork, userID string) ([]map[string]interface{}, error) {
	if m.getBlacklistFunc != nil {
		return m.getBlacklistFunc(socialNetwork, userID)
	}
	return nil, nil
}

func (m *MockDatabase) CheckBlacklist(socialNetwork, userID, nickname string) (bool, error) {
	if m.checkBlacklistFunc != nil {
		return m.checkBlacklistFunc(socialNetwork, userID, nickname)
	}
	return false, nil
}

// TestHandlePostBlacklistSuccess тестирует успешное добавление в чёрный список
func TestHandlePostBlacklistSuccess(t *testing.T) {
	mockDB := &MockDatabase{
		addToBlacklistFunc: func(socialNetwork, userID, nickname string) (int, error) {
			return 1, nil
		},
	}

	handler := HandlePostBlacklist(mockDB)

	payload := map[string]string{
		"social_network": "telegram",
		"user_id":        "550e8400-e29b-41d4-a716-446655440000",
		"nickname":       "baduser",
	}
	body, _ := json.Marshal(payload)

	req := httptest.NewRequest("POST", "/api/post", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, w.Code)
	}

	var resp ResponsePayload
	json.NewDecoder(w.Body).Decode(&resp)

	if resp.Status != "success" {
		t.Errorf("Expected status 'success', got '%s'", resp.Status)
	}
}

// TestHandlePostBlacklistInvalidMethod тестирует некорректный HTTP метод
func TestHandlePostBlacklistInvalidMethod(t *testing.T) {
	mockDB := &MockDatabase{}
	handler := HandlePostBlacklist(mockDB)

	req := httptest.NewRequest("GET", "/api/post", nil)
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	if w.Code != http.StatusMethodNotAllowed {
		t.Errorf("Expected status %d, got %d", http.StatusMethodNotAllowed, w.Code)
	}
}

// TestHandlePostBlacklistInvalidJSON тестирует некорректный JSON
func TestHandlePostBlacklistInvalidJSON(t *testing.T) {
	mockDB := &MockDatabase{}
	handler := HandlePostBlacklist(mockDB)

	req := httptest.NewRequest("POST", "/api/post", bytes.NewReader([]byte("invalid json")))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status %d, got %d", http.StatusBadRequest, w.Code)
	}
}

// TestHandlePostBlacklistInvalidSocialNetwork тестирует некорректную социальную сеть
func TestHandlePostBlacklistInvalidSocialNetwork(t *testing.T) {
	mockDB := &MockDatabase{}
	handler := HandlePostBlacklist(mockDB)

	payload := map[string]string{
		"social_network": "facebook",
		"user_id":        "550e8400-e29b-41d4-a716-446655440000",
		"nickname":       "baduser",
	}
	body, _ := json.Marshal(payload)

	req := httptest.NewRequest("POST", "/api/post", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status %d, got %d", http.StatusBadRequest, w.Code)
	}

	var resp ResponsePayload
	json.NewDecoder(w.Body).Decode(&resp)

	if resp.Message != "Invalid social network" {
		t.Errorf("Expected message 'Invalid social network', got '%s'", resp.Message)
	}
}

// TestHandlePostBlacklistMissingUserID тестирует отсутствие user_id
func TestHandlePostBlacklistMissingUserID(t *testing.T) {
	mockDB := &MockDatabase{}
	handler := HandlePostBlacklist(mockDB)

	payload := map[string]string{
		"social_network": "telegram",
		"user_id":        "",
		"nickname":       "baduser",
	}
	body, _ := json.Marshal(payload)

	req := httptest.NewRequest("POST", "/api/post", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status %d, got %d", http.StatusBadRequest, w.Code)
	}

	var resp ResponsePayload
	json.NewDecoder(w.Body).Decode(&resp)

	if resp.Message != "User ID is required" {
		t.Errorf("Expected message 'User ID is required', got '%s'", resp.Message)
	}
}

// TestHandlePostBlacklistMissingNickname тестирует отсутствие nickname
func TestHandlePostBlacklistMissingNickname(t *testing.T) {
	mockDB := &MockDatabase{}
	handler := HandlePostBlacklist(mockDB)

	payload := map[string]string{
		"social_network": "telegram",
		"user_id":        "550e8400-e29b-41d4-a716-446655440000",
		"nickname":       "",
	}
	body, _ := json.Marshal(payload)

	req := httptest.NewRequest("POST", "/api/post", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status %d, got %d", http.StatusBadRequest, w.Code)
	}

	var resp ResponsePayload
	json.NewDecoder(w.Body).Decode(&resp)

	if resp.Message != "Nickname is required" {
		t.Errorf("Expected message 'Nickname is required', got '%s'", resp.Message)
	}
}

// TestHandleGetBlacklistSuccess тестирует успешное получение чёрного списка
func TestHandleGetBlacklistSuccess(t *testing.T) {
	mockDB := &MockDatabase{
		getBlacklistFunc: func(socialNetwork, userID string) ([]map[string]interface{}, error) {
			return []map[string]interface{}{
				{
					"nickname":   "baduser1",
					"created_at": "2025-01-01T00:00:00Z",
				},
				{
					"nickname":   "baduser2",
					"created_at": "2025-01-02T00:00:00Z",
				},
			}, nil
		},
	}

	handler := HandleGetBlacklist(mockDB)

	req := httptest.NewRequest("GET", "/api/blacklist?social_network=telegram&user_id=123", nil)
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, w.Code)
	}

	var resp ResponsePayload
	json.NewDecoder(w.Body).Decode(&resp)

	if resp.Status != "success" {
		t.Errorf("Expected status 'success', got '%s'", resp.Status)
	}

	data := resp.Data.(map[string]interface{})
	if data["banned_count"] != float64(2) {
		t.Errorf("Expected banned_count 2, got %v", data["banned_count"])
	}
}

// TestHandleGetBlacklistInvalidMethod тестирует некорректный HTTP метод
func TestHandleGetBlacklistInvalidMethod(t *testing.T) {
	mockDB := &MockDatabase{}
	handler := HandleGetBlacklist(mockDB)

	req := httptest.NewRequest("POST", "/api/blacklist", nil)
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	if w.Code != http.StatusMethodNotAllowed {
		t.Errorf("Expected status %d, got %d", http.StatusMethodNotAllowed, w.Code)
	}
}

// TestHandleGetBlacklistMissingSocialNetwork тестирует отсутствие social_network
func TestHandleGetBlacklistMissingSocialNetwork(t *testing.T) {
	mockDB := &MockDatabase{}
	handler := HandleGetBlacklist(mockDB)

	req := httptest.NewRequest("GET", "/api/blacklist?user_id=550e8400-e29b-41d4-a716-446655440000", nil)
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status %d, got %d", http.StatusBadRequest, w.Code)
	}

	var resp ResponsePayload
	json.NewDecoder(w.Body).Decode(&resp)

	if resp.Message != "social_network parameter is required" {
		t.Errorf("Expected message 'social_network parameter is required', got '%s'", resp.Message)
	}
}

// TestHandleGetBlacklistMissingUserID тестирует отсутствие user_id
func TestHandleGetBlacklistMissingUserID(t *testing.T) {
	mockDB := &MockDatabase{}
	handler := HandleGetBlacklist(mockDB)

	req := httptest.NewRequest("GET", "/api/blacklist?social_network=telegram", nil)
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status %d, got %d", http.StatusBadRequest, w.Code)
	}

	var resp ResponsePayload
	json.NewDecoder(w.Body).Decode(&resp)

	if resp.Message != "user_id parameter is required" {
		t.Errorf("Expected message 'user_id parameter is required', got '%s'", resp.Message)
	}
}

// TestHandleGetBlacklistInvalidSocialNetwork тестирует некорректную социальную сеть
func TestHandleGetBlacklistInvalidSocialNetwork(t *testing.T) {
	mockDB := &MockDatabase{}
	handler := HandleGetBlacklist(mockDB)

	req := httptest.NewRequest("GET", "/api/blacklist?social_network=facebook&user_id=550e8400-e29b-41d4-a716-446655440000", nil)
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status %d, got %d", http.StatusBadRequest, w.Code)
	}

	var resp ResponsePayload
	json.NewDecoder(w.Body).Decode(&resp)

	if resp.Message != "Invalid social network" {
		t.Errorf("Expected message 'Invalid social network', got '%s'", resp.Message)
	}
}

// TestHandleCheckBlacklistSuccess тестирует успешную проверку (пользователь в списке)
func TestHandleCheckBlacklistSuccess(t *testing.T) {
	mockDB := &MockDatabase{
		checkBlacklistFunc: func(socialNetwork, userID, nickname string) (bool, error) {
			return true, nil
		},
	}

	handler := HandleCheckBlacklist(mockDB)

	req := httptest.NewRequest("GET", "/api/check?social_network=telegram&user_id=550e8400-e29b-41d4-a716-446655440000&nickname=baduser", nil)
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, w.Code)
	}

	var resp ResponsePayload
	json.NewDecoder(w.Body).Decode(&resp)

	if resp.Status != "success" {
		t.Errorf("Expected status 'success', got '%s'", resp.Status)
	}

	data := resp.Data.(map[string]interface{})
	if data["is_banned"] != true {
		t.Errorf("Expected is_banned true, got %v", data["is_banned"])
	}
}

// TestHandleCheckBlacklistUserNotInList тестирует проверку (пользователь не в списке)
func TestHandleCheckBlacklistUserNotInList(t *testing.T) {
	mockDB := &MockDatabase{
		checkBlacklistFunc: func(socialNetwork, userID, nickname string) (bool, error) {
			return false, nil
		},
	}

	handler := HandleCheckBlacklist(mockDB)

	req := httptest.NewRequest("GET", "/api/check?social_network=telegram&user_id=550e8400-e29b-41d4-a716-446655440000&nickname=gooduser", nil)
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, w.Code)
	}

	var resp ResponsePayload
	json.NewDecoder(w.Body).Decode(&resp)

	data := resp.Data.(map[string]interface{})
	if data["is_banned"] != false {
		t.Errorf("Expected is_banned false, got %v", data["is_banned"])
	}
}

// TestHandleCheckBlacklistInvalidMethod тестирует некорректный HTTP метод
func TestHandleCheckBlacklistInvalidMethod(t *testing.T) {
	mockDB := &MockDatabase{}
	handler := HandleCheckBlacklist(mockDB)

	req := httptest.NewRequest("POST", "/api/check", nil)
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	if w.Code != http.StatusMethodNotAllowed {
		t.Errorf("Expected status %d, got %d", http.StatusMethodNotAllowed, w.Code)
	}
}

// TestHandleCheckBlacklistMissingSocialNetwork тестирует отсутствие social_network
func TestHandleCheckBlacklistMissingSocialNetwork(t *testing.T) {
	mockDB := &MockDatabase{}
	handler := HandleCheckBlacklist(mockDB)

	req := httptest.NewRequest("GET", "/api/check?user_id=550e8400-e29b-41d4-a716-446655440000&nickname=baduser", nil)
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status %d, got %d", http.StatusBadRequest, w.Code)
	}

	var resp ResponsePayload
	json.NewDecoder(w.Body).Decode(&resp)

	if resp.Message != "social_network parameter is required" {
		t.Errorf("Expected message 'social_network parameter is required', got '%s'", resp.Message)
	}
}

// TestHandleCheckBlacklistMissingUserID тестирует отсутствие user_id
func TestHandleCheckBlacklistMissingUserID(t *testing.T) {
	mockDB := &MockDatabase{}
	handler := HandleCheckBlacklist(mockDB)

	req := httptest.NewRequest("GET", "/api/check?social_network=telegram&nickname=baduser", nil)
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status %d, got %d", http.StatusBadRequest, w.Code)
	}

	var resp ResponsePayload
	json.NewDecoder(w.Body).Decode(&resp)

	if resp.Message != "user_id parameter is required" {
		t.Errorf("Expected message 'user_id parameter is required', got '%s'", resp.Message)
	}
}

// TestHandleCheckBlacklistMissingNickname тестирует отсутствие nickname
func TestHandleCheckBlacklistMissingNickname(t *testing.T) {
	mockDB := &MockDatabase{}
	handler := HandleCheckBlacklist(mockDB)

	req := httptest.NewRequest("GET", "/api/check?social_network=telegram&user_id=550e8400-e29b-41d4-a716-446655440000", nil)
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status %d, got %d", http.StatusBadRequest, w.Code)
	}

	var resp ResponsePayload
	json.NewDecoder(w.Body).Decode(&resp)

	if resp.Message != "nickname parameter is required" {
		t.Errorf("Expected message 'nickname parameter is required', got '%s'", resp.Message)
	}
}

// TestHandleCheckBlacklistInvalidSocialNetwork тестирует некорректную социальную сеть
func TestHandleCheckBlacklistInvalidSocialNetwork(t *testing.T) {
	mockDB := &MockDatabase{}
	handler := HandleCheckBlacklist(mockDB)

	req := httptest.NewRequest("GET", "/api/check?social_network=facebook&user_id=550e8400-e29b-41d4-a716-446655440000&nickname=baduser", nil)
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status %d, got %d", http.StatusBadRequest, w.Code)
	}

	var resp ResponsePayload
	json.NewDecoder(w.Body).Decode(&resp)

	if resp.Message != "Invalid social network" {
		t.Errorf("Expected message 'Invalid social network', got '%s'", resp.Message)
	}
}

// TestHandleRoot тестирует корневой endpoint
func TestHandleRoot(t *testing.T) {
	req := httptest.NewRequest("GET", "/", nil)
	w := httptest.NewRecorder()

	HandleRoot(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, w.Code)
	}

	var resp map[string]string
	json.NewDecoder(w.Body).Decode(&resp)

	if resp["message"] != "Welcome to Warden Solver Core" {
		t.Errorf("Expected message 'Welcome to Warden Solver Core', got '%s'", resp["message"])
	}

	if resp["version"] != "1.0.0" {
		t.Errorf("Expected version '1.0.0', got '%s'", resp["version"])
	}
}

// BenchmarkHandlePostBlacklist бенчмарк для POST запроса
func BenchmarkHandlePostBlacklist(b *testing.B) {
	mockDB := &MockDatabase{
		addToBlacklistFunc: func(socialNetwork, userID, nickname string) (int, error) {
			return 1, nil
		},
	}

	handler := HandlePostBlacklist(mockDB)
	payload := map[string]string{
		"social_network": "telegram",
		"user_id":        "550e8400-e29b-41d4-a716-446655440000",
		"nickname":       "baduser",
	}
	body, _ := json.Marshal(payload)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		req := httptest.NewRequest("POST", "/api/post", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		handler.ServeHTTP(w, req)
	}
}

// BenchmarkHandleGetBlacklist бенчмарк для GET запроса списка
func BenchmarkHandleGetBlacklist(b *testing.B) {
	mockDB := &MockDatabase{
		getBlacklistFunc: func(socialNetwork, userID string) ([]map[string]interface{}, error) {
			return []map[string]interface{}{
				{"nickname": "baduser1"},
			}, nil
		},
	}

	handler := HandleGetBlacklist(mockDB)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		req := httptest.NewRequest("GET", "/api/blacklist?social_network=telegram&user_id=550e8400-e29b-41d4-a716-446655440000", nil)
		w := httptest.NewRecorder()
		handler.ServeHTTP(w, req)
	}
}

// BenchmarkHandleCheckBlacklist бенчмарк для проверки в списке
func BenchmarkHandleCheckBlacklist(b *testing.B) {
	mockDB := &MockDatabase{
		checkBlacklistFunc: func(socialNetwork, userID, nickname string) (bool, error) {
			return true, nil
		},
	}

	handler := HandleCheckBlacklist(mockDB)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		req := httptest.NewRequest("GET", "/api/check?social_network=telegram&user_id=550e8400-e29b-41d4-a716-446655440000&nickname=baduser", nil)
		w := httptest.NewRecorder()
		handler.ServeHTTP(w, req)
	}
}
