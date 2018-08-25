Simple RESTful API to manage book collections. This project was written using the Go and MongoDB as a database.

## Getting Started

These instructions will get you a copy of the project up and running on your local machine for development and testing purposes.

### Prerequisities

- Go >= 1.10.3
- MongoDB

### Installing

1. Clone this repo

2. Fetching Dependecies

   [TOML](https://github.com/BurntSushi/toml) parser and encoder
   
   ```
   go get github.com/BurntSushi/toml
   ```
   
   Package [mux](https://github.com/gorilla/mux) implements a request router and dispatcher
    
   ```
   go get github.com/gorilla/mux
   ```

   [MongoDB driver](https://github.com/go-mgo/mgo)
   
   ```
   go get gopkg.in/mgo.v2
   ```

3. Run `go run app.go`

