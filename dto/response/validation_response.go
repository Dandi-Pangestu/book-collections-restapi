package response

type ValidationResponse struct {
	Error string `json:"error"`
	Message string `json:"message"`
	Errors []string `json:"errors"`
}
