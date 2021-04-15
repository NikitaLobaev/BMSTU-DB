package response

type Response struct {
	Code int
	JSONObject interface{}
}

func NewResponse(code int, JSONObject interface{}) *Response {
	return &Response{
		Code: code,
		JSONObject: JSONObject,
	}
}
