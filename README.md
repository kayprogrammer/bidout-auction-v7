# BIDOUT AUCTION V7
A simple bidding API built with Fiber 

![alt text](https://github.com/kayprogrammer/bidout-auction-v7/blob/main/display/fiber.png?raw=true)


#### FIBER DOCS: [Documentation](https://docs.gofiber.io/)
#### GORM DOCS: [Documentation](https://gorm.io/)
#### PG ADMIN: [Documentation](https://pgadmin.org) 


## How to run locally

* Download this repo or run: 
```bash
    $ git clone git@github.com:kayprogrammer/bidout-auction-v7.git
```

#### In the root directory:
- Install all dependencies
```bash
    $ go install github.com/cosmtrek/air@latest 
    $ go mod download
```
- Create an `.env` file and copy the contents from the `.env.example` to the file and set the respective values. A postgres database can be created with PG ADMIN or psql

- Run Locally
```bash
    $ air
```

- Run With Docker
```bash
    $ docker-compose up --build -d --remove-orphans
```
OR
```bash
    $ make build
```

- Test Coverage
```bash
    $ go test ./tests -v -count=1
```
OR
```bash
    $ make test
```

## Docs
#### API Url: [BidOut Docs](http://127.0.0.1:8000/) 
#### Swagger: [Documentation](https://swagger.io/docs/)

![alt text](https://github.com/kayprogrammer/bidout-auction-v7/blob/main/display/display1.png?raw=true)

![alt text](https://github.com/kayprogrammer/bidout-auction-v7/blob/main/display/display2.png?raw=true)

![alt text](https://github.com/kayprogrammer/bidout-auction-v7/blob/main/display/display3.png?raw=true)

![alt text](https://github.com/kayprogrammer/bidout-auction-v7/blob/main/display/display4.png?raw=true)

![alt text](https://github.com/kayprogrammer/bidout-auction-v7/blob/main/display/display5.png?raw=true)

![alt text](https://github.com/kayprogrammer/bidout-auction-v7/blob/main/display/display6.png?raw=true)

![alt text](https://github.com/kayprogrammer/bidout-auction-v7/blob/main/display/display7.png?raw=true)