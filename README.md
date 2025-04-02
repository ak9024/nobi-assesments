# nobi-assessment

API service for NOBI Investment products using Golang Fiber and MySQL.

## Endpoint List Changes  

I have made some changes to the endpoint list to make it neater and easier to read.

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

Just for notes this project just covering one cycle transaction with assumption NAB value keep with default value == 1, if want to add to adjust NAB, just hit endpoints `/api/investments` with POST method to create new investment with adjustable NAB value.

### API Endpoints

Customers

POST /api/customers - Create a new customer

GET /api/customers - Get a list of customers

GET /api/customers/{customer_uuid} - Get details of a specific customer

Investments

POST /api/investments - Create a new investment

GET /api/investments - Get a list of investments

GET /api/investments/{investment_uuid} - Get details of a specific investment

Transactions

POST /api/transactions/deposit - Make a deposit transaction

POST /api/transactions/withdraw - Make a withdrawal transaction

GET /api/transactions/customer/{customer_uuid} - Get transactions for a specific customer

Portfolio

GET /api/portfolio/{customer_id}/{investment_id} - Get portfolio details for a customer and investment

### For testing

If you have `docker` just run `sh run.sh` to run docker compose and integration testing for the api, please check ./tests/api_test.go for references.

```shell
# clone the repository
git clone https://github.com/username/nobi-investment.git
cd nobi-investment
# run all project dependencies
docker compose up -d
```

Then open postman and import the collection from `./nobi-assesment.postman_collection.json`

### For development

```shell
# clone the repository
git clone https://github.com/username/nobi-investment.git
cd nobi-investment

# download dependencies
go mod tidy
# copy env for configuration
cp .env.example .env
# run mysql
docker compose up -d mysql
go run .
```

For more references please check `./nobi-assesment.postman_collection.json` postman collection for the API.
