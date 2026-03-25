package messages

type Response struct {
	WinnerAmount int
}

func NewResponse(WinnerAmount int) Response {
	response := Response{WinnerAmount}
	return response
}

func (r Response) Type() byte {
	return MSG_RESPONSE
}
