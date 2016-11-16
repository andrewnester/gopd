package gopd

import (
	"net/http"
	"encoding/json"
)

type WebHookData struct {
	Document
	Event      string `json:"event"`
	ActionDate string `json:"action_date"`
	ActionBy   User `json:"action_by"`
}

func (w *WebHookData) FromRequest(r *http.Request) (*WebHookData, error) {
	decoder := json.NewDecoder(r.Body)
	var wh WebHookData
	err := decoder.Decode(&wh)
	if err != nil {
		return nil, err
	}
	defer r.Body.Close()
	return &wh, nil
}
