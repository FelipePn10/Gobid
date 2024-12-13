package services

import (
	"context"
	"errors"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"log/slog"
	"sync"
)

type MessageKind = int

const (
	PlaceBid MessageKind = iota
	SuccessFullyPlaceBid
	FailedToPlaceBid
	NewBidPlaced
	AuctionFinished
)

type Message struct {
	Message string
	Amount  float64
	Kind    MessageKind
	UserId  uuid.UUID
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

func (r *AuctionRoom) registerClient(c *Client) {
	slog.Info("New user Connected", "Client", c)
	r.Clients[c.UserId] = c
}

func (r *AuctionRoom) unregisterClient(c *Client) {
	slog.Info("User disconnected", "Client", c)
	delete(r.Clients, c.UserId)
	r.Clients[c.UserId] = c
}

func (r *AuctionRoom) broadcastMessage(m Message) {
	slog.Info("New message recieved", "RoomID", r.Id, "message", m.Message, "user_id", m.UserId)
	switch m.Kind {
	case PlaceBid:
		bid, err := r.BidsService.Placebid(r.Context, r.Id, m.UserId, m.Amount)
		if err != nil {
			if errors.Is(err, ErrBidIsTooLow) {
				if client, ok := r.Clients[m.UserId]; ok {
					client.Send <- Message{Kind: FailedToPlaceBid, Message: ErrBidIsTooLow.Error()}
				}
				return
			}
		}
		if client, ok := r.Clients[m.UserId]; ok {
			client.Send <- Message{Kind: SuccessFullyPlaceBid, Message: "Your bid was successfully placed."}
		}
		for id, client := range r.Clients {
			newBidMessage := Message{Kind: NewBidPlaced, Message: "A new bid was placed.", Amount: bid.BidAmount}
			if id == m.UserId {
				continue
			}
			client.Send <- newBidMessage
		}
	}
}
func (r *AuctionRoom) Run() {

	defer func() {
		close(r.Broadcast)
		close(r.Register)
		close(r.Unregister)
	}()
	for {
		select {
		case client := <-r.Register:

			r.registerClient(client)
		case client := <-r.Unregister:

			r.unregisterClient(client)
		case message := <-r.Broadcast:
			r.broadcastMessage(message)
		case <-r.Context.Done():
			slog.Info("Auction has ended.", "auctionID", r.Id)
			for _, client := range r.Clients {
				client.Send <- Message{Kind: AuctionFinished, Message: "auction has been finished"}
			}
			return
		}
	}
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
