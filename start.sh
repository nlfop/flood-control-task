trap "docker-compose stop" EXIT
go mod tidy
docker-compose up -d
sleep 5
cd flood_control/
go test
cd ..
export NFLOOD=3
export KFLOOD=4
go run main.go

