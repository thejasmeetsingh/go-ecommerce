package api

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"
	"github.com/thejasmeetsingh/go-ecommerce/order_service/internal/database"
	"github.com/thejasmeetsingh/go-ecommerce/order_service/shared"
)

func getUserID(c *gin.Context) (uuid.UUID, error) {
	userIDStr, isExists := c.Get("userID")
	if !isExists {
		return uuid.Nil, fmt.Errorf("authentication failed")
	}

	userID, err := uuid.Parse(fmt.Sprintf("%v", userIDStr))
	if err != nil {
		return uuid.Nil, err
	}

	return userID, nil
}

// API for creating an order record
func (apiCfg *APIConfig) CreateOrder(c *gin.Context) {
	type Parameters struct {
		ProductID string `json:"product_id" binding:"required"`
	}
	var params Parameters

	err := c.ShouldBindJSON(&params)
	if err != nil {
		log.Error("Error caught while parsing order creation request data: ", err)
		c.JSON(http.StatusBadRequest, gin.H{"message": "Error while parsing the request"})
		return
	}

	userID, err := getUserID(c)
	if err != nil {
		log.Error("Error caught while parsing user id: ", err)
		c.JSON(http.StatusForbidden, gin.H{"message": "Invalid credentails"})
		return
	}

	_, err = uuid.Parse(params.ProductID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid product id format"})
		return
	}

	product, err := shared.GetProductIDToDetails(apiCfg.Cache, c, params.ProductID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid product id"})
		return
	}

	order, err := CreateOrderDB(apiCfg, c, database.CreateOrderParams{
		ID:         uuid.New(),
		CreatedAt:  time.Now().UTC(),
		ModifiedAt: time.Now().UTC(),
		UserID:     userID,
		ProductID:  product.ID,
	}, product)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": order})
}

// API for getting order list
func (apiCfg *APIConfig) GetOrders(c *gin.Context) {
	// Parse string offset to integer
	offsetStr := c.Query("offset")
	if offsetStr == "" {
		offsetStr = "0"
	}

	offset, err := strconv.Atoi(offsetStr)
	if err != nil {
		offset = 0
	}

	userID, err := getUserID(c)
	if err != nil {
		log.Error("Error caught while parsing user id: ", err)
		c.JSON(http.StatusForbidden, gin.H{"message": "Invalid credentails"})
		return
	}

	orders, err := GetOrderListDB(apiCfg, c, database.GetOrdersParams{
		UserID: userID,
		Limit:  10,
		Offset: int32(offset),
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	if len(orders) == 0 {
		c.JSON(http.StatusOK, gin.H{"results": []string{}})
		return
	}

	c.JSON(http.StatusOK, gin.H{"results": orders})
}

// API for getting details of a specific order by its ID
func (apiCfg *APIConfig) GetOrderDetail(c *gin.Context) {
	orderIDStr := c.Param("id")
	orderID, err := uuid.Parse(orderIDStr)

	if err != nil {
		log.Error("Error caught while parsing the order ID: ", err)
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid order id format"})
		return
	}

	userID, err := getUserID(c)
	if err != nil {
		log.Error("Error caught while parsing user id: ", err)
		c.JSON(http.StatusForbidden, gin.H{"message": "Invalid credentails"})
		return
	}

	order, err := GetOrderDetailDB(apiCfg, c, database.GetOrderByIdParams{
		ID:     orderID,
		UserID: userID,
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": order})
}

// API for deleting an order by its ID
func (apiCfg *APIConfig) DeleteOrder(c *gin.Context) {
	orderIDStr := c.Param("id")
	orderID, err := uuid.Parse(orderIDStr)

	if err != nil {
		log.Error("Error caught while parsing the order ID: ", err)
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid order id format"})
		return
	}

	userID, err := getUserID(c)
	if err != nil {
		log.Error("Error caught while parsing user id: ", err)
		c.JSON(http.StatusForbidden, gin.H{"message": "Invalid credentails"})
		return
	}

	_, err = GetOrderDetailDB(apiCfg, c, database.GetOrderByIdParams{
		ID:     orderID,
		UserID: userID,
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	err = DeleteOrderDB(apiCfg, c, database.DeleteOrderParams{
		ID:     orderID,
		UserID: userID,
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Order deleted successfully!"})
}
