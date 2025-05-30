package response

import (
	"encoding/json"
	"net/http"
	"strconv"

	// "github.com/bytedance/sonic"

	"github.com/wagecloud/wagecloud-server/internal/logger"
	"go.uber.org/zap"
)

type Error struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

type CommonResponse struct {
	Data  any    `json:"data,omitempty"`
	Error *Error `json:"error,omitempty"`
}

type PaginateResponse[T any] struct {
	Data       []T        `json:"data"`
	Pagination Pagination `json:"pagination"`
	Error      *Error     `json:"errors,omitempty"`
}

type Pagination struct {
	Limit      int32   `json:"limit"`
	Page       int32   `json:"page"`
	Total      int64   `json:"total"`
	NextPage   *int32  `json:"next_page"`
	NextCursor *string `json:"next_cursor"`
}

func writeError(w http.ResponseWriter, errCode string, message string) {
	response, err := json.Marshal(CommonResponse{
		Data: nil,
		Error: &Error{
			Code:    errCode,
			Message: message,
		},
	})
	if err != nil {
		w.Write([]byte("Error marshalling JSON"))
		return
	}
	w.Write(response)
}

func writeResponse(w http.ResponseWriter, dto any) {
	response, err := json.Marshal(dto)
	if err != nil {
		writeError(w, http.StatusText(http.StatusInternalServerError), "Error marshalling JSON")
		return
	}

	w.Write(response)
}

func FromDTO(w http.ResponseWriter, dto any, httpCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(httpCode)

	writeResponse(w, CommonResponse{
		Data:  dto,
		Error: nil,
	})
}

func FromError(w http.ResponseWriter, err error, httpCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(httpCode)

	// // Internal server error
	// var errWithCode *model.ErrorWithCode
	// if errors.As(err, &errWithCode) {
	// 	logger.Log.Error("error", zap.Error(errWithCode.Err))
	// 	writeError(w, errWithCode.Code, errWithCode.Msg)
	// 	return
	// }

	// Normal http error
	logger.Log.Error("error", zap.Error(err))
	writeError(w, strconv.Itoa(httpCode), err.Error())
}

func FromHTTPError(w http.ResponseWriter, httpCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(httpCode)

	writeError(w, strconv.Itoa(httpCode), http.StatusText(httpCode))
}

func FromPaginate[T any](w http.ResponseWriter, paginateResult PaginateResult[T]) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	writeResponse(w, PaginateResponse[T]{
		Data: paginateResult.Data,
		Pagination: Pagination{
			Limit:      paginateResult.Limit,
			Page:       paginateResult.Page,
			Total:      paginateResult.Total,
			NextPage:   paginateResult.NextPage,
			NextCursor: paginateResult.NextCursor,
		},
		Error: nil,
	})
}

// PaginateResult represents a paginated result set
type PaginateResult[T any] struct {
	Data       []T     `json:"data"`
	Limit      int32   `json:"limit"`
	Page       int32   `json:"page"`
	Total      int64   `json:"total"`
	NextPage   *int32  `json:"next_page,omitempty"`
	NextCursor *string `json:"next_cursor,omitempty"`
}
