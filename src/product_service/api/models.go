package api

import (
	"time"

	"github.com/google/uuid"
	"github.com/thejasmeetsingh/go-ecommerce/product_service/internal/database"
)

type productDetail struct {
	ID          uuid.UUID `json:"id"`
	CreatedAt   time.Time `json:"created_at"`
	ModifiedAt  time.Time `json:"modified_at"`
	Name        string    `json:"name"`
	Price       int32     `json:"price"`
	Description string    `json:"description"`
}

type productList struct {
	ID          uuid.UUID `json:"id"`
	Name        string    `json:"name"`
	Price       int32     `json:"price"`
	Description string    `json:"description"`
}

func DatabaseProductToProduct(dbProduct database.Product) productDetail {
	return productDetail{
		ID:          dbProduct.ID,
		CreatedAt:   dbProduct.CreatedAt,
		ModifiedAt:  dbProduct.ModifiedAt,
		Name:        dbProduct.Name,
		Price:       dbProduct.Price,
		Description: dbProduct.Description,
	}
}

func DatabaseProductToProductList(dbProducts []database.GetProductsRow) []productList {
	var products []productList

	for _, dbProduct := range dbProducts {
		products = append(products, productList{
			ID:          dbProduct.ID,
			Name:        dbProduct.Name,
			Price:       dbProduct.Price,
			Description: dbProduct.Description,
		})
	}
	return products
}
