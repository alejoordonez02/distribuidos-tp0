package messages

type Response struct {
	Winners []string
}

func NewResponse(Winners []string) Response {
	response := Response{Winners}
	return response
}

func (r Response) Type() byte {
	return MSG_RESPONSE
}
