package gamesocket

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/NiklasPrograms/tictacgo/backend/pkg/game"
)

type Command interface {
	execute() ([]ResponseHandler, error)
	// parseContent parses the content of the command and returns an error if the content is invalid. The content is then stored in the command as a field.
	parseContent(any) error
}

func NewCommand(message GameMessage, client *GameClient) (Command, error) {
	command, err := parseInstruction(message.Instruction, client)
	if err != nil {
		return nil, err
	}

	if err := command.parseContent(message.Content); err != nil {
		return nil, err
	}

	return command, nil
}

func parseInstruction(instruction string, client *GameClient) (Command, error) {
	instruction = strings.TrimSpace(strings.ToLower(instruction))

	switch instruction {
	case "start game":
		return NewStartGameCommand(client)
	case "choose square":
		return NewChooseSquareCommand(client)
	case "select character":
		return NewSelectCharacterCommand(client)
	case "request draw":
		return NewRequestDrawCommand(client)
	case "respond to draw":
		return NewRespondToDrawRequestCommand(client)
	case "withdraw draw request":
		return NewWithdrawDrawRequestCommand(client)
	default:
		return nil, fmt.Errorf("invalid game instruction: %s", instruction)
	}
}

// ------------------------------------------------------------------------------------------------------------

type StartGameCommand struct {
	client *GameClient
}

func NewStartGameCommand(client *GameClient) (*StartGameCommand, error) {
	if client == nil {
		return nil, fmt.Errorf("Client cannot be nil")
	}
	return &StartGameCommand{client}, nil
}

func (*StartGameCommand) parseContent(any) error {
	return nil // StartGameCommand has no content
}

func (command *StartGameCommand) execute() ([]ResponseHandler, error) {
	pool := command.client.Pool
	game := pool.game

	var response GameResponse

	bothCharactersSelected := pool.xClient != nil && pool.oClient != nil
	if !bothCharactersSelected { // TODO kunne overveje at lave en ErrorResponse, der bare sendes tilbage til den pågældende client og ikke broadcastes
		return nil, fmt.Errorf("Both characters must be selected, before a game can start")
	}

	isClientPlaying := pool.xClient == command.client || pool.oClient == command.client
	if !isClientPlaying {
		return nil, fmt.Errorf("Client must be playing to start the game, cannot be a spectator")
	}

	game.StartGame()

	response.ResponseType = GAME_STARTED

	responseHandler, err := NewResponseHandler(&response, pool.clients)
	if err != nil {
		return nil, err
	}

	return []ResponseHandler{*responseHandler}, nil
}

// ------------------------------------------------------------------------------------------------------------

type ChooseSquareCommand struct {
	client   *GameClient
	position game.Position
}

func NewChooseSquareCommand(client *GameClient) (*ChooseSquareCommand, error) {
	if client == nil {
		return nil, fmt.Errorf("Client cannot be nil")
	}

	command := &ChooseSquareCommand{
		client:   client,
		position: game.NO_POSITION,
	}

	return command, nil

}

func (command *ChooseSquareCommand) parseContent(content any) error {
	position, err := game.ParsePosition(content)
	if err != nil {
		return err
	}
	command.position = position
	return nil
}

func (command *ChooseSquareCommand) isClientInTurn() bool {
	pool := command.client.Pool

	if pool.game.PlayerInTurn() == game.X {
		return pool.xClient == command.client
	}
	return pool.oClient == command.client
}

func (command *ChooseSquareCommand) execute() ([]ResponseHandler, error) {
	pool := command.client.Pool

	var response GameResponse

	if command.position == game.NO_POSITION {
		return nil, fmt.Errorf("Position cannot be empty")
	}

	if !command.isClientInTurn() {
		return nil, fmt.Errorf("It was not this client's turn to play")
	}

	board, err := pool.game.ChooseSquare(command.position)
	if err != nil {
		return nil, err
	}

	response.ResponseType = BOARD
	response.Body = board

	responseHandler, err := NewResponseHandler(&response, pool.clients)
	if err != nil {
		return nil, err
	}

	return []ResponseHandler{*responseHandler}, nil
}

// ------------------------------------------------------------------------------------------------------------

type SelectCharacterCommand struct {
	client    *GameClient
	character game.SquareCharacter
}

func NewSelectCharacterCommand(client *GameClient) (*SelectCharacterCommand, error) {
	if client == nil {
		return nil, fmt.Errorf("Client cannot be nil")
	}

	command := &SelectCharacterCommand{
		client:    client,
		character: game.EMPTY_CHARACTER,
	}

	return command, nil
}

func (command *SelectCharacterCommand) parseContent(content any) error {
	character, err := game.ParseSquareCharacter(content)
	if err != nil {
		return err
	}
	command.character = character
	return nil
}

func (command *SelectCharacterCommand) selectCharacter() error {
	pool := command.client.Pool

	if command.character == game.X && pool.xClient == nil {
		pool.xClient = command.client
		return nil
	}

	if command.character == game.O && pool.oClient == nil {
		pool.oClient = command.client
		return nil
	}

	return fmt.Errorf("Character %v is already taken", command.character)
}

func (command *SelectCharacterCommand) execute() ([]ResponseHandler, error) {
	pool := command.client.Pool

	var response GameResponse

	if command.character == game.EMPTY_CHARACTER {
		return nil, fmt.Errorf("Character cannot be empty")
	}

	hasClientAlreadySelected := pool.xClient == command.client || pool.oClient == command.client
	if hasClientAlreadySelected {
		return nil, fmt.Errorf("Client had already selected a character")
	}

	if err := command.selectCharacter(); err != nil {
		return nil, err
	}

	response.ResponseType = CHARACTER_SELECTED
	response.Body = command.character

	responseHandler, err := NewResponseHandler(&response, pool.clients)
	if err != nil {
		return nil, err
	}

	return []ResponseHandler{*responseHandler}, nil
}

// ------------------------------------------------------------------------------------------------------------

type RequestDrawCommand struct {
	client *GameClient
}

func NewRequestDrawCommand(client *GameClient) (*RequestDrawCommand, error) {
	if client == nil {
		return nil, fmt.Errorf("Client cannot be nil")
	}
	return &RequestDrawCommand{client}, nil
}

func (*RequestDrawCommand) parseContent(any) error {
	return nil // RequestDrawCommand has no content
}

func (command *RequestDrawCommand) execute() ([]ResponseHandler, error) {
	pool := command.client.Pool

	if pool.DrawRequestHandler.IsDrawRequested {
		return nil, fmt.Errorf("A draw has already been requested")
	}

	isValid, err := IsDrawRequestCommandValid(command.client)
	if !isValid {
		return nil, err
	}

	responseHandlers, err := CreateResponseHandlersForSelfAndOpponent(command.client, REQUEST_DRAW)
	if err != nil {
		return nil, err
	}

	pool.DrawRequestHandler.IsDrawRequested = true
	pool.DrawRequestHandler.DrawRequester = command.client

	return responseHandlers, nil
}

// ------------------------------------------------------------------------------------------------------------

type RespondToDrawRequestCommand struct {
	client *GameClient
	accept bool
}

func NewRespondToDrawRequestCommand(client *GameClient) (*RespondToDrawRequestCommand, error) {
	if client == nil {
		return nil, fmt.Errorf("Client cannot be nil")
	}

	command := &RespondToDrawRequestCommand{
		client: client,
		accept: false,
	}

	return command, nil
}

func (command *RespondToDrawRequestCommand) parseContent(content any) error {
	accept, ok := content.(bool)
	if !ok {
		return fmt.Errorf("invalid content type for RespondToDrawInstruction: %v", reflect.TypeOf(content))
	}
	command.accept = accept
	return nil
}

// TODO should add a time limit to respond to a draw request, probably through an attached timestamp?
func (command *RespondToDrawRequestCommand) execute() ([]ResponseHandler, error) {
	pool := command.client.Pool
	game := pool.game

	var response GameResponse
	var receivers []*GameClient

	if !pool.DrawRequestHandler.IsDrawRequested {
		return nil, fmt.Errorf("No draw was requested, so there is nothing to respond to")
	}

	isValid, err := IsDrawRequestCommandValid(command.client)
	if !isValid {
		return nil, err
	}

	isClientTheRequester := pool.DrawRequestHandler.DrawRequester == command.client
	if isClientTheRequester {
		return nil, fmt.Errorf("Client cannot respond to their own draw request")
	}

	if command.accept {
		game.ForceDraw()
		receivers = pool.clients
	} else {
		receivers = []*GameClient{pool.xClient, pool.oClient} // Only opponent and self needs to know that the draw was rejected. No spectators receive the inital message either
	}

	pool.DrawRequestHandler.IsDrawRequested = false
	pool.DrawRequestHandler.DrawRequester = nil

	response.Body = command.accept
	response.ResponseType = DRAW_REQUEST_RESPONSE

	responseHandler, err := NewResponseHandler(&response, receivers)
	if err != nil {
		return nil, err
	}

	return []ResponseHandler{*responseHandler}, nil
}

// ------------------------------------------------------------------------------------------------------------

type WithdrawDrawRequestCommand struct {
	client *GameClient
}

func NewWithdrawDrawRequestCommand(client *GameClient) (*WithdrawDrawRequestCommand, error) {
	if client == nil {
		return nil, fmt.Errorf("Client cannot be nil")
	}

	return &WithdrawDrawRequestCommand{client}, nil
}

func (*WithdrawDrawRequestCommand) parseContent(any) error {
	return nil // WithdrawDrawRequestCommand has no content
}

func (command *WithdrawDrawRequestCommand) execute() ([]ResponseHandler, error) {
	pool := command.client.Pool

	if !pool.DrawRequestHandler.IsDrawRequested {
		return nil, fmt.Errorf("No draw was requested, so there is nothing to withdraw")
	}

	isValid, err := IsDrawRequestCommandValid(command.client)
	if !isValid {
		return nil, err
	}

	isClientTheDrawRequester := pool.DrawRequestHandler.DrawRequester == command.client
	if !isClientTheDrawRequester {
		return nil, fmt.Errorf("Client is not the draw requester, and cannot withdraw the draw request")
	}

	responseHandlers, err := CreateResponseHandlersForSelfAndOpponent(command.client, WITHDRAW_DRAW_REQUEST)
	if err != nil {
		return nil, err
	}

	pool.DrawRequestHandler.IsDrawRequested = false
	pool.DrawRequestHandler.DrawRequester = nil

	return responseHandlers, nil
}
