package v1

// ResponsePayload представляет структуру ответа
type ResponsePayload struct {
	Status  string      `json:"status"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

// DatabaseInterface определяет интерфейс для работы с БД
type DatabaseInterface interface {
	AddToBlacklist(socialNetwork, userID, nickname string) (int, error)
	GetBlacklist(socialNetwork, userID string) ([]map[string]interface{}, error)
	CheckBlacklist(socialNetwork, userID, nickname string) (bool, error)
}

// SocialNetworks определяет допустимые социальные сети
var SocialNetworks = map[string]bool{
	"x":        true,
	"telegram": true,
	"twitch":   true,
}
