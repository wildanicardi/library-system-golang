package helper

import (
	"encoding/json"
	"net/http"
)

// response json
func RenderJSON(res http.ResponseWriter, statusCode int, data interface{}) {
	res.WriteHeader(statusCode)
	res.Header().Set("Content-Type", "application/json")
	json.NewEncoder(res).Encode(data)
}
