package responses

func NewErrorResponse(m string) Response {
	return Response{
		Success: false,
		Data:    nil,
		Message: m,
	}
}

func NewErrorResponseWithData(m string, d interface{}) Response {
	return Response{
		Success: false,
		Data:    d,
		Message: m,
	}
}
