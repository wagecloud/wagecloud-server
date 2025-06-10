package response

type Error struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

type ErrorWithCode interface {
	Error() string
	Code() string
}
