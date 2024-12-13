package services

import (
	"context"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"sync"
)

type MessageKind = int

const (
	PlaceBid MessageKind = iota
)

type Message struct {
	Message string
	Kind    MessageKind
}

type AuctionLobby struct {
	sync.Mutex
	Rooms map[uuid.UUID]*AuctionRoom
}

type AuctionRoom struct {
	Id          uuid.UUID
	Context     context.Context
	Broadcast   chan Message
	Register    chan *Client
	Unregister  chan *Client
	Clients     map[uuid.UUID]*Client
	BidsService *BidsService
}

type Client struct {
	Room   *AuctionRoom
	Conn   *websocket.Conn
	Send   chan Message
	UserId uuid.UUID
}

func NewAuctionRoom(ctx context.Context, id uuid.UUID, BidsService BidsService) *AuctionRoom {
	return &AuctionRoom{
		Id:          id,
		Context:     ctx,
		Broadcast:   make(chan Message),
		Register:    make(chan *Client),
		Unregister:  make(chan *Client),
		Clients:     make(map[uuid.UUID]*Client),
		BidsService: &BidsService,
	}
}

func NewClient(room *AuctionRoom, userId uuid.UUID, conn *websocket.Conn) *Client {
	return &Client{
		Room:   room,
		Conn:   conn,
		UserId: userId,
		Send:   make(chan Message, 512),
	}
}
