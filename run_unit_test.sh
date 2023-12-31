#!/bin/bash/

container_name=test_db
db_user=db_user
db_password=1234
db_name=test

services=("user_service" "product_service")

# Run the testing database container
docker run --name $container_name \
        -e POSTGRES_USER=$db_user \
        -e POSTGRES_PASSWORD=$db_password \
        -e POSTGRES_DB=$db_name \
        -p 5432:5432 \
        --health-cmd='pg_isready -d $db_name -U $db_user' \
        --health-interval=10s \
        --health-timeout=5s \
        --health-retries=5 \
        -d postgres:16.1-alpine3.18

echo "Waiting for DB..."

while true; do
    # Check the health of the database container
    if docker inspect --format '{{json .State.Health.Status}}' $container_name | grep -q "healthy"
    then
        break
    fi
    sleep 1
done

sleep 30

for service in ${services[@]};
do
    echo "Runnig test for $service service"

    cd src/$service
    goose -dir sql/schema postgres postgres://$db_user:$db_password@localhost:5432/$db_name up
    go test ./api/

    echo "---------------------"
done

# Stop & remove docker container
docker container stop $container_name
docker container rm $container_name
docker volume prune -a -f