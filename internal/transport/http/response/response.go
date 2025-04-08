package response

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	// "github.com/bytedance/sonic"

	"github.com/wagecloud/wagecloud-server/internal/logger"
	"github.com/wagecloud/wagecloud-server/internal/model"
)

func writeError(w http.ResponseWriter, errCode string, message string) {
	response, err := json.Marshal(CommonResponse{
		Message: message,
		Data:    nil,
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

func FromDTO(w http.ResponseWriter, dto any, httpCode int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(httpCode)

	writeResponse(w, CommonResponse{
		Message: message,
		Data:    dto,
		Error:   nil,
	})
}

func FromError(w http.ResponseWriter, err error) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusInternalServerError)

	// Service error
	var errWithCode *model.ErrorWithCode
	if errors.As(err, &errWithCode) {
		writeError(w, errWithCode.Code, errWithCode.Msg)
		return
	}

	// Internal server error
	// writeError(w, http.StatusText(http.StatusInternalServerError), err.Error())
	logger.Log.Error(err.Error())
	writeError(w, http.StatusText(http.StatusInternalServerError), "Internal Server Error")
}

func FromHTTPError(w http.ResponseWriter, httpCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(httpCode)

	writeError(w, strconv.Itoa(httpCode), http.StatusText(httpCode))
}

func FromPaginate[T any](w http.ResponseWriter, paginateResult model.PaginateResult[T], message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	writeResponse(w, PaginateResponse[T]{
		Message: message,
		Data:    paginateResult.Data,
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
