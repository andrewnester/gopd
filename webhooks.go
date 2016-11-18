package gopd

import (
	"net/http"
	"encoding/json"
)

type WebHook struct {
	Event      string `json:"event"`
	Data WebHookData `json:"data"`
}

type WebHookData struct {
	Document
	ActionDate string `json:"action_date"`
	ActionBy   User `json:"action_by"`
}

func (w *WebHookData) FromRequest(r *http.Request) ([]WebHook, error) {
	decoder := json.NewDecoder(r.Body)
	var whs []WebHook
	err := decoder.Decode(&whs)
	if err != nil {
		return nil, err
	}
	defer r.Body.Close()
	return whs, nil
}
