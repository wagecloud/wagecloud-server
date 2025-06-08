package response

type CommonResponse struct {
	Data  any    `json:"data,omitempty"`
	Error *Error `json:"error,omitempty"`
}
