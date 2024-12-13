package services

import (
	"context"
	"database/sql"
	"errors"
	"github.com/FelipePn10/Gobid/internal/store/pgstore"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"time"
)

type ProductsService struct {
	pool    *pgxpool.Pool
	queries *pgstore.Queries
}

var ErrProductNotFound = errors.New("product not found")

func NewProductsService(pool *pgxpool.Pool) ProductsService {
	return ProductsService{
		pool:    pool,
		queries: pgstore.New(pool),
	}
}

func (ps *ProductsService) CreateProduct(
	ctx context.Context,
	sellerId uuid.UUID,
	productName,
	description string,
	baseprice float64,
	auctionEnd time.Time,
) (uuid.UUID, error) {
	id, err := ps.queries.CreatedProduct(ctx, pgstore.CreatedProductParams{
		SellerID:    sellerId,
		ProductName: productName,
		Description: description,
		Baseprice:   baseprice,
		AuctionEnd:  auctionEnd,
	})
	if err != nil {
		return uuid.UUID{}, err
	}
	return id, nil
}

func (ps *ProductsService) UpdateProduct(
	ctx context.Context,
	productID uuid.UUID,
	sellerID uuid.UUID,
	productName *string,
	description *string,
	basePrice *float64,
	auctionEnd *time.Time,
) error {
	params := pgstore.UpdateProductParams{
		ID:          productID,
		SellerID:    sellerID,
		ProductName: nullString(productName),
		Description: nullString(description),
		Baseprice:   nullFloat64(basePrice),
		AuctionEnd:  nullTime(auctionEnd),
	}
	if err := ps.queries.UpdateProduct(ctx, params); err != nil {
		return err
	}

	return nil
}

func nullString(s *string) sql.NullString {
	if s == nil {
		return sql.NullString{Valid: false}
	}
	return sql.NullString{String: *s, Valid: true}
}

func nullFloat64(f *float64) sql.NullFloat64 {
	if f == nil {
		return sql.NullFloat64{Valid: false}
	}
	return sql.NullFloat64{Float64: *f, Valid: true}
}

func nullTime(t *time.Time) sql.NullTime {
	if t == nil {
		return sql.NullTime{Valid: false}
	}
	return sql.NullTime{Time: *t, Valid: true}
}

func (ps *ProductsService) DeleteProduct(
	ctx context.Context,
	productID uuid.UUID,
	sellerID uuid.UUID,
) error {
	err := ps.queries.DeleteProduct(ctx, pgstore.DeleteProductParams{
		ID:       productID,
		SellerID: sellerID,
	})
	if err != nil {
		return err
	}
	return nil
}

func (ps *ProductsService) GetProductByID(ctx context.Context, productID uuid.UUID) (pgstore.Product, error) {
	product, err := ps.queries.GetProductById(ctx, productID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return pgstore.Product{}, ErrProductNotFound
		}
		return pgstore.Product{}, err
	}
	return product, nil
}
