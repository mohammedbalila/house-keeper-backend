## House Keeper API

### About the project
- An experimental project to try go [echo](https://github.com/labstack/echo) framework.
- It was meant to solve a small issue we're having at home.
- The code is not well organised and missing some production level enhancements.
- All PRs are welcome!. Maybe it'll be suitable for prod in the future? XD.
### The issue it solves
We're a group of friends sharing the same house hold and having troubles keeping up with our payments as a group. Sometimes one of us would pay a bill by himself and collect the money later from the others, but when there are too many bills to keep up with, things get ugly real quick, so It's like a reminder for people to keep up with their bills and pay those who handles paying a bill by themselves.
### The stack
- Go [echo](https://github.com/labstack/echo) v4 API
- PostgreSQL >= 14
- [Flyway](https://flywaydb.org/) for schema management
- [Flutter](https://flutter.dev/) front-end
### Run the project locally
 ```bash
 $ go mod download
 # get the packages

 $ go run main.go
 # Run the server
 ```

### Development
-	[Go](https://go.dev/doc/install) >= 1.17
-	[PostgreSQL](https://www.postgresql.org/download/) >= 14
-	Go [air](https://github.com/cosmtrek/air) 

### Schema
-	[Database schema](https://github.com/mustafabalila/house-keeper-db-schema)
 ### Environment

 See [config](./.env.sample).
