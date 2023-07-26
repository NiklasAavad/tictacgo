package gamesocket

type ResponseHandler struct {
	Response  *GameResponse
	Receivers []*GameClient
}

func NewResponseHandler(response *GameResponse, receivers []*GameClient) (*ResponseHandler, error) {
	// assume for now that the response is valid and the receivers are valid
	return &ResponseHandler{response, receivers}, nil
}
