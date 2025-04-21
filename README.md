
# Flaco

Git of a GO project using GRPC in the RT0805 module of the Master DAS in Reims

## Authors

- Corentin
- Flavien


## Usage



```bash

# Go to the "Code" folder
cd Code/

go mod tidy

# Run docker-compose
docker-compose up --build

# Run main file
go run main.go
```

## Unit Tests



```bash
# Go to the "Code" folder
cd Code/

# Run client test file
go test client/client_test.go

# Run server test file
go test serveur/serveur_test.go
```
