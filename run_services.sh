#!/bin/bash/

# Run User service
docker-compose -f src/user_service/docker-compose.yml up -d
echo "User services are up and running!"

# Run Product service
docker-compose -f src/product_service/docker-compose.yml up -d
echo "Product services are up and running!"

# Run Order service
docker-compose -f src/order_service/docker-compose.yml up -d
echo "Order services are up and running!"