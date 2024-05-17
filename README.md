
# Flaco

Git of a GO project using GRPC in the RT0805 module of the Master DAS in Reims

## Authors

- Dupont Corentin
- Morlet Flavien


## Usage



```bash

# Go to the "Code" folder:
cd Code/

go mod tidy

# Run docker-compose:
docker-compose up --build

go run main.go
```

## Unit Tests



```bash
Go to the "Code" folder:
cd Code/

go test tests/client_test.go
go test tests/serveur_test.go
```