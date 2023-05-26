package utils

// Response 响应结果结构体定义
type Response struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Result  interface{} `json:"result"`
}

func Success(result interface{}) *Response {
	return &Response{
		Success: true,
		Message: "成功",
		Result:  result,
	}
}

func SuccessWithMessage(result interface{}, message string) *Response {
	return &Response{
		Success: true,
		Message: message,
		Result:  result,
	}
}

func Fail(message string) *Response {
	return &Response{
		Success: false,
		Message: message,
		Result:  nil,
	}
}
