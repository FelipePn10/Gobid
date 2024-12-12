package api

import (
	"github.com/FelipePn10/Gobid/internal/jsonutils"
	"github.com/FelipePn10/Gobid/internal/usecase/product"
	"github.com/google/uuid"
	"net/http"
)

func (api *Api) handleCreateProduct(w http.ResponseWriter, r *http.Request) {
	data, problems, err := jsonutils.DecodeValidJson[product.CreateProductReq](r)
	if err != nil {
		jsonutils.EncodeJson(w, r, http.StatusUnprocessableEntity, problems)
		return
	}
	userID, ok := api.Sessions.Get(r.Context(), "AuthenticatedUserId").(uuid.UUID)
	if !ok {
		jsonutils.EncodeJson(w, r, http.StatusInternalServerError, map[string]any{
			"error": "unexpected error, try again later",
		})
		return
	}
	id, err := api.ProductService.CreateProduct(
		r.Context(),
		userID,
		data.ProductName,
		data.Description,
		data.Baseprice,
		data.AuctionEnd,
	)
	if err != nil {
		jsonutils.EncodeJson(w, r, http.StatusInternalServerError, map[string]any{
			"error": "failed to create product, try again later",
		})
	}
	jsonutils.EncodeJson(w, r, http.StatusInternalServerError, map[string]any{
		"message":    "product created successfully",
		"product_id": id,
	})
}
