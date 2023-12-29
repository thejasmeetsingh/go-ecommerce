#!/bin/bash/

# Run the testing database container
docker run --name test_users_db \
        -e POSTGRES_USER=db_user \
        -e POSTGRES_PASSWORD=1234 \
        -e POSTGRES_DB=users_test_db \
        -p 5432:5432 \
        --health-cmd='pg_isready -d users_test_db -U db_user' \
        --health-interval=10s \
        --health-timeout=5s \
        --health-retries=5 \
        -d postgres:16.1-alpine3.18

echo "Waiting for DB..."

while true; do
    # Check the health of the database container
    if docker inspect --format '{{json .State.Health.Status}}' test_users_db | grep -q "healthy"
    then
        break
    fi
    sleep 1
done

sleep 30

# Run migrations
goose -dir sql/schema postgres postgres://db_user:1234@localhost:5432/users_test_db up

# Run test case
go test ./api/

# Stop & remove docker container
docker container stop test_users_db
docker container rm test_users_db
docker volume prune -a -f
