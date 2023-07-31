package gamesocket

import "errors"

type ResponseHandler struct {
	Response  *GameResponse
	Receivers []*GameClient
}

func NewResponseHandler(response *GameResponse, receivers []*GameClient) (*ResponseHandler, error) {
	// assume for now that the response is valid and the receivers are valid
	return &ResponseHandler{response, receivers}, nil
}

// Creates a response handler for the client and the opponent
// The body of the response is a boolean indicating whether the receiver is the opponent
func CreateResponseHandlersForSelfAndOpponent(self *GameClient, responseType ResponseType) ([]ResponseHandler, error) {
	pool := self.Pool

	isBothCharactersChosen := pool.xClient != nil && pool.oClient != nil
	if !isBothCharactersChosen {
		return nil, errors.New("Cannot send a response to the opponent if both characters are not chosen")
	}

	isOpponent := true
	responseToOpponent := GameResponse{
		ResponseType: responseType,
		Body:         isOpponent,
	}

	isOpponent = false
	responseToSelf := GameResponse{
		ResponseType: responseType,
		Body:         isOpponent,
	}

	var opponent *GameClient
	if pool.xClient == self {
		opponent = pool.oClient
	} else {
		opponent = pool.xClient
	}

	opponentResponseHandler, err := NewResponseHandler(&responseToOpponent, []*GameClient{opponent})
	if err != nil {
		return nil, err
	}

	selfResponseHandler, err := NewResponseHandler(&responseToSelf, []*GameClient{self})
	if err != nil {
		return nil, err
	}

	return []ResponseHandler{*opponentResponseHandler, *selfResponseHandler}, nil
}
