package responses

func NewResourceResponse(d interface{}) Response {
	return Response{
		Success: true,
		Data:    d,
		Message: "success",
	}
}

func MessageResponse(m string) Response {
	return Response{
		Success: true,
		Data:    nil,
		Message: m,
	}
}

func ResourceResponseWithMessage(m string, d interface{}) Response {
	return Response{
		Success: true,
		Data:    d,
		Message: m,
	}
}
