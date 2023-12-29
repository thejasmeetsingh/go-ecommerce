#!/bin/bash/

# Run User service
docker-compose -f src/user_service/docker-compose.yml up -d

# Run Product service
docker-compose -f src/product_service/docker-compose.yml up -d

# RUn Order service
docker-compose -f src/order_service/docker-compose.yml up -d