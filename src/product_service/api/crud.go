package api

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	log "github.com/sirupsen/logrus"
	"github.com/thejasmeetsingh/go-ecommerce/product_service/internal/database"
)

func structToJson(product database.Product) ([]byte, error) {
	return json.Marshal(product)
}

func JsonToStruct(data []byte) (database.Product, error) {
	var product database.Product
	err := json.Unmarshal(data, &product)
	if err != nil {
		return product, err
	}
	return product, nil
}

func retreiveProductFromCache(client *redis.Client, ctx *gin.Context, key string) ([]byte, error) {
	return client.Get(ctx, key).Bytes()
}

func storeProductToCache(client *redis.Client, ctx *gin.Context, product database.Product) error {
	key := product.ID.String()
	value, err := structToJson(product)
	if err != nil {
		return err
	}
	return client.Set(ctx, key, value, 1*time.Hour).Err()
}

// Create product
func CreateProductDB(apiCfg *APIConfig, ctx *gin.Context, params database.CreateProductParams) (database.Product, error) {
	// Begin DB transaction
	tx, err := apiCfg.DB.Begin()
	if err != nil {
		log.Fatal(err)
		return database.Product{}, fmt.Errorf("something went wrong")
	}
	defer tx.Rollback()
	qtx := apiCfg.Queries.WithTx(tx)

	dbProduct, err := qtx.CreateProduct(ctx, params)

	if err != nil {
		log.Errorln(err)
		return database.Product{}, fmt.Errorf("something went wrong")
	}

	// Store newly product details into cache
	err = storeProductToCache(apiCfg.Cache, ctx, dbProduct)
	if err != nil {
		log.Error(err)
		return database.Product{}, fmt.Errorf("something went wrong")
	}

	// Commit the transaction
	err = tx.Commit()
	if err != nil {
		log.Fatal(err)
		return database.Product{}, fmt.Errorf("something went wrong")
	}

	return dbProduct, nil
}

// Get list of products
func GetProductListDB(apiCfg *APIConfig, ctx *gin.Context, params database.GetProductsParams) ([]database.GetProductsRow, error) {
	products, err := apiCfg.Queries.GetProducts(ctx, params)
	if err != nil {
		return []database.GetProductsRow{}, fmt.Errorf("something went wrong")
	}
	return products, nil
}

// Get details of a specific product
func GetProductDetailDB(apiCfg *APIConfig, ctx *gin.Context, productID uuid.UUID) (database.Product, error) {
	// Find if product is available in cache or not
	data, err := retreiveProductFromCache(apiCfg.Cache, ctx, productID.String())
	if err != nil {
		product, err := apiCfg.Queries.GetProductById(ctx, productID)
		if err != nil {
			log.Errorln(err)
			return database.Product{}, fmt.Errorf("something went wrong")
		}

		// store the product details into cache
		err = storeProductToCache(apiCfg.Cache, ctx, product)
		if err != nil {
			log.Error(err)
			return database.Product{}, fmt.Errorf("something went wrong")
		}

		return product, nil
	}

	// Convert product JSON to struct
	product, err := JsonToStruct(data)
	if err != nil {
		log.Errorln(err)
		return database.Product{}, fmt.Errorf("something went wrong")
	}

	return product, nil
}

// Update product details
func UpdateProductDetailDB(apiCfg *APIConfig, ctx *gin.Context, params database.UpdateProductDetailsParams) (database.Product, error) {
	// Begin DB transaction
	tx, err := apiCfg.DB.Begin()
	if err != nil {
		log.Fatal(err)
		return database.Product{}, fmt.Errorf("something went wrong")
	}
	defer tx.Rollback()
	qtx := apiCfg.Queries.WithTx(tx)

	product, err := qtx.UpdateProductDetails(ctx, params)

	if err != nil {
		log.Errorln(err)
		return database.Product{}, fmt.Errorf("something went wrong")
	}

	// Store product with latest details to cache
	err = storeProductToCache(apiCfg.Cache, ctx, product)
	if err != nil {
		log.Error(err)
		return database.Product{}, fmt.Errorf("something went wrong")
	}

	// Commit the transaction
	err = tx.Commit()
	if err != nil {
		log.Fatal(err)
		return database.Product{}, fmt.Errorf("something went wrong")
	}

	return product, nil
}

// Delete a product
func DeleteProductDetailDB(apiCfg *APIConfig, ctx *gin.Context, productID uuid.UUID) error {
	// Begin DB transaction
	tx, err := apiCfg.DB.Begin()
	if err != nil {
		log.Fatal(err)
		return fmt.Errorf("something went wrong")
	}
	defer tx.Rollback()
	qtx := apiCfg.Queries.WithTx(tx)

	err = qtx.DeleteProduct(ctx, productID)
	if err != nil {
		log.Errorln(err)
		return fmt.Errorf("something went wrong")
	}

	// Remove deleted product from the cache if it is preasent
	err = apiCfg.Cache.Del(ctx, productID.String()).Err()
	if err != nil {
		log.Error(err)
	}

	// Commit the transaction
	err = tx.Commit()
	if err != nil {
		log.Fatal(err)
		return fmt.Errorf("something went wrong")
	}

	return nil
}
