# Misinfo Detector Backend

This is the backend for the misinformation detector service for COS30049. It sits
in front of a SQLite database to store user posts. It has the 
[Pymisinfo AI engine](https://github.com/codenor/pymisinfo-monorepo) running in the 
background, and will use it to determine the validity of certain posts.

## Running The Program

### Database Connection Strings

If you are using this program to test it, you may just use the in-memory database
with `:memory:` and ignore the rest of this section. Otheriwse, please follow the 
[Go Sqlite driver documentation](https://github.com/mattn/go-sqlite3?tab=readme-ov-file#connection-string)
for information about connection strings, and how to configure one best for your 
use case.

### Locally

If you have Go installed, you may simply run the program using the command line.

```
go run .
```

There are some command line arguments you may use. You can pass the `-h` flag to
view them. 

```
Usage of misinfodetector-backend:
  -listen string
    	where this program should listen for api requests (default 127.0.0.1:5000)
  -sqlite string
    	where the sqlite database should be stored (default :memory:)
```

### Docker

If you don't like using the really useful and handy Go cli tools, you may use the 
boring Docker compose script. 

Before running, you must configure your environment variables by copying the examples.
Please view this file. There are sensible defaults, but its good to know about the 
things you can change in case you have problems.

```
cp .env.exmaple .env
```

You may run the following command to run the container:

```
docker-compose up --build
```

## Testing

This project has unit/integration tests configured. You will need to have this project
working locally for this to work. Run the following command to execute all tests:

```
go test ./...
```

See `go help test` for more details for how to run these tests.
