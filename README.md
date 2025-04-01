# nobi-assessment

API service for NOBI Investment products using Golang Fiber and MySQL.

### Structures

```
❯ tree -L 2
.
├── README.md
├── cmd
│   └── api
├── database.sql
├── delivery
│   └── http
├── docker-compose.yaml
├── go.mod
├── go.sum
├── internal
│   ├── domain
│   ├── repository
│   └── usecase
├── nobi-assesment-v2.postman_collection.json
└── pkg
    ├── db
    └── utils

12 directories, 6 files
```

### Tech Stack

- Golang - Main programming language
- Fiber Framework - Fast and efficient web framework for Golang
- MySQL - Relational database for storing data
- UUID - For unique identification of each entity (replacing auto increment)
- Docker (optional) - For containerization and deployment


### Getting Started

```shell
# clone the repository
git clone https://github.com/username/nobi-investment.git
cd nobi-investment

# download dependencies
go mod tidy
# copy env for configuration
cp .env.example .env
# run mysql
docker compose up -d
```

For more references please check `./nobi-assesment.postman_collection.json` postman collection for the API.
