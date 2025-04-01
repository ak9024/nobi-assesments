docker compose up -d --build
sleep 15
go test ./tests/api_test.go
docker compose down -v
