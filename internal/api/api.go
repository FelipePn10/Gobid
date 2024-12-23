package api

import (
	"github.com/FelipePn10/Gobid/internal/services"
	"github.com/alexedwards/scs/v2"
	"github.com/go-chi/chi/v5"
	"github.com/gorilla/websocket"
)

type Api struct {
	Router         *chi.Mux
	UserService    services.UserService
	ProductService services.ProductsService
	Sessions       *scs.SessionManager
	WsUpgrader     websocket.Upgrader
	AuctionLobby   services.AuctionLobby
	BidsService    services.BidsService
}
