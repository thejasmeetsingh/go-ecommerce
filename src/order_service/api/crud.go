package api

import (
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"github.com/thejasmeetsingh/go-ecommerce/order_service/internal/database"
	"github.com/thejasmeetsingh/go-ecommerce/order_service/models"
	"github.com/thejasmeetsingh/go-ecommerce/order_service/shared"
)

func CreateOrderDB(apiCfg *APIConfig, ctx *gin.Context, params database.CreateOrderParams, product models.Product) (models.OrderDetail, error) {
	// Begin DB transaction
	tx, err := apiCfg.DB.Begin()
	if err != nil {
		log.Fatal("Error caught while initiating the transaction: ", err)
		return models.OrderDetail{}, fmt.Errorf("something went wrong")
	}
	defer tx.Rollback()
	qtx := apiCfg.Queries.WithTx(tx)

	// Create the order
	dbOrder, err := qtx.CreateOrder(ctx, params)

	if err != nil {
		log.Errorln("Error caught while creating the order in DB", err)
		return models.OrderDetail{}, fmt.Errorf("something went wrong")
	}

	order := models.DatabaseOrderToOrder(dbOrder, product)

	// Convert the order details to bytes
	orderByte, err := models.OrderStructToByte(order)
	if err != nil {
		log.Errorln("Error caught while converting order struct obj to byte: ", err)
		return models.OrderDetail{}, fmt.Errorf("something went wrong")
	}

	// Store newly product details into cache
	err = apiCfg.Cache.Set(ctx, order.ID.String(), orderByte, 1*time.Hour).Err()
	if err != nil {
		log.Error("Error caught while saving order details to cache: ", err)
		return models.OrderDetail{}, fmt.Errorf("something went wrong")
	}

	// Commit the transaction
	err = tx.Commit()
	if err != nil {
		log.Fatal("Error caught while committing the transaction; ", err)
		return models.OrderDetail{}, fmt.Errorf("something went wrong")
	}

	return order, nil
}

func GetOrderListDB(apiCfg *APIConfig, ctx *gin.Context, params database.GetOrdersParams) ([]models.OrderList, error) {
	dbOrders, err := apiCfg.Queries.GetOrders(ctx, params)
	if err != nil {
		log.Error("Error caught while fetching order list: ", err)
		return []models.OrderList{}, fmt.Errorf("something went wrong")
	}
	orders := models.DatabaseOrderToOrderList(dbOrders)
	return orders, nil
}

func GetOrderDetailDB(apiCfg *APIConfig, ctx *gin.Context, params database.GetOrderByIdParams) (models.OrderDetail, error) {
	orderByte, err := apiCfg.Cache.Get(ctx, params.ID.String()).Bytes()
	if err != nil {
		dbOrder, err := apiCfg.Queries.GetOrderById(ctx, params)
		if err != nil {
			log.Error("Error caught while fetching order details: ", err)
			return models.OrderDetail{}, fmt.Errorf("something went wrong")
		}

		// Fetch product details
		product, err := shared.GetProductIDToDetails(apiCfg.Cache, ctx, dbOrder.ProductID.String())
		if err != nil {
			return models.OrderDetail{}, fmt.Errorf("something went wrong")
		}

		// Convert the order details to bytes
		order := models.OrderDetail{
			ID:         dbOrder.ID,
			CreatedAt:  dbOrder.CreatedAt,
			ModifiedAt: dbOrder.ModifiedAt,
			Product:    product,
		}
		orderByte, err := models.OrderStructToByte(order)
		if err != nil {
			log.Error("Error caught while converting order details to byte: ", err)
			return models.OrderDetail{}, fmt.Errorf("something went wrong")
		}

		// Store order details in cache
		err = apiCfg.Cache.Set(ctx, params.ID.String(), orderByte, 1*time.Hour).Err()
		if err != nil {
			log.Error("Error caught while saving order details to cache: ", err)
			return models.OrderDetail{}, fmt.Errorf("something went wrong")
		}

		return order, nil
	}

	// Convert cached order bytes to order model
	order, err := models.ByteToOrderStruct(orderByte)
	if err != nil {
		log.Error("Error caught while converting order byte to order details: ", err)
		return models.OrderDetail{}, fmt.Errorf("something went wrong")
	}

	return order, nil
}

func DeleteOrderDB(apiCfg *APIConfig, ctx *gin.Context, params database.DeleteOrderParams) error {
	// Begin DB transaction
	tx, err := apiCfg.DB.Begin()
	if err != nil {
		log.Fatal("Error caught while initiating the transaction: ", err)
		return fmt.Errorf("something went wrong")
	}
	defer tx.Rollback()
	qtx := apiCfg.Queries.WithTx(tx)

	// Delete the order
	err = qtx.DeleteOrder(ctx, params)
	if err != nil {
		log.Error("Error caught while deleting order details: ", err)
		return fmt.Errorf("something went wrong")
	}

	// Remove deleted product from the cache
	err = apiCfg.Cache.Del(ctx, params.ID.String()).Err()
	if err != nil {
		log.Error("Error caught while deleting the order details from cache: ", err)
	}

	// Commit the transaction
	err = tx.Commit()
	if err != nil {
		log.Fatal("Error caught while committing the transaction: ", err)
		return fmt.Errorf("something went wrong")
	}

	return nil
}
