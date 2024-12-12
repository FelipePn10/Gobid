package product

import (
	"context"
	"github.com/FelipePn10/Gobid/internal/validator"
	"github.com/google/uuid"
	"time"
)

type CreateProductReq struct {
	SellerID    uuid.UUID `json:"seller_id"`
	ProductName string    `json:"product_name"`
	Description string    `json:"description"`
	Baseprice   float64   `json:"baseprice"`
	AuctionEnd  time.Time `json:"auction_end"`
}

const minAuctionDuration = 2 * time.Hour

func (req CreateProductReq) Valid(ctx context.Context) validator.Evaluator {
	var eval validator.Evaluator

	eval.CheckField(validator.NotBlank(req.ProductName), "product_name", "this field cannot be blank")
	eval.CheckField(validator.NotBlank(req.Description), "description", "this field cannot be blank")
	eval.CheckField(validator.MinChars(req.Description, 35) &&
		validator.MaxChars(req.Description, 3500), "description", "your description must have a minimum of 35 and a maximum of 3500 characters")
	eval.CheckField(req.Baseprice > 0, "baseprice", "the product value must be at least greater than zero")
	eval.CheckField(req.AuctionEnd.Sub(time.Now()) >= minAuctionDuration, "auction_end", "the duration must be at least two hours")
	return eval
}
