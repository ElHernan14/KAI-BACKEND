package response

type APIResponse struct {
	Status       string      `json:"status"`
	Code         int         `json:"code"`
	ErrorMessage string      `json:"errorMessage,omitempty"`
	Data         interface{} `json:"data,omitempty"`
}

func Success(data interface{}) APIResponse {
	return APIResponse{
		Status: "success",
		Code:   200,
		Data:   data,
	}
}

func SuccessWithCode(code int, data interface{}) APIResponse {
	return APIResponse{
		Status: "success",
		Code:   code,
		Data:   data,
	}
}

func Error(code int, msg string) APIResponse {
	return APIResponse{
		Status:       "error",
		Code:         code,
		ErrorMessage: msg,
		Data:         nil,
	}
}
