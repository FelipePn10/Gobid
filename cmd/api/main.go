package main

import (
	"context"
	"encoding/gob"
	"fmt"
	"github.com/FelipePn10/Gobid/internal/api"
	"github.com/FelipePn10/Gobid/internal/services"
	"github.com/alexedwards/scs/pgxstore"
	"github.com/alexedwards/scs/v2"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
	"net/http"
	"os"
	"time"
)

func main() {
	gob.Register(uuid.UUID{})

	if err := godotenv.Load(); err != nil {
		panic(err)
	}
	ctx := context.Background()
	pool, err := pgxpool.New(ctx, fmt.Sprintf("user=%s password=%s host=%s port=%s dbname=%s sslmode=disable",
		os.Getenv("GOBID_DATABASE_USER"),
		os.Getenv("GOBID_DATABASE_PASSWORD"),
		os.Getenv("GOBID_DATABASE_HOST"),
		os.Getenv("GOBID_DATABASE_PORT"),
		os.Getenv("GOBID_DATABASE_NAME"),
	))
	if err != nil {
		panic(err)
	}

	defer pool.Close()

	if err := pool.Ping(ctx); err != nil {
		panic(err)
	}

	s := scs.New()
	s.Store = pgxstore.New(pool)
	s.Lifetime = 24 * time.Hour
	s.Cookie.HttpOnly = true
	s.Cookie.SameSite = http.SameSiteLaxMode

	api := api.Api{
		Router:         chi.NewMux(),
		UserService:    services.NewUserService(pool),
		ProductService: services.NewProductsService(pool),
		BidsService:    services.NewBidsService(pool),
		Sessions:       s,
		WsUpgrader: websocket.Upgrader{
			CheckOrigin: func(r *http.Request) bool {
				allowedOrigins := []string{
					"http://localhost:3080",
				}
				origin := r.Header.Get("Origin")
				for _, allowedOrigin := range allowedOrigins {
					if origin == allowedOrigin {
						return true
					}
				}
				return false
			},
		},
		AuctionLobby: services.AuctionLobby{
			Rooms: make(map[uuid.UUID]*services.AuctionRoom),
		},
	}
	api.BindRoutes()

	fmt.Println("Starting server on port :3080")
	if err := http.ListenAndServe("localhost:3080", api.Router); err != nil {
		panic(err)
	}
}
