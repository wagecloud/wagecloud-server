package response

import (
	"encoding/json"
	"net/http"

	// "github.com/bytedance/sonic"
	"github.com/wagecloud/wagecloud-server/internal/model"
)

func FromDTO(w http.ResponseWriter, code int, dto any) {
	response, err := json.Marshal(dto)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Error marshalling JSON"))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}

func FromDTOWithMessage(w http.ResponseWriter, dto any, code int, message string) {
	response, err := json.Marshal(dto)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Error marshalling JSON"))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}

func FromMessage(w http.ResponseWriter, code int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write([]byte(`{"message": "` + message + `"}`))
}

func FromError(w http.ResponseWriter, code int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write([]byte(`{"message": "` + message + `"}`))
}

func FromPaginate[T any](w http.ResponseWriter, paginateResult model.PaginateResult[T]) {
	response, err := json.Marshal(paginateResult)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Error marshalling JSON"))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(response)
}
