<p align="center">
    <img src="https://raw.githubusercontent.com/jedsaxon/misinfodetector-backend/refs/heads/master/docs/troll-cage.png" alt="trolls in a cage thanks to chad golang gopher" />
</p>

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

You can also configure environment variables. Any command line arguments you use will
override everything in `.env`, so be careful when doing this. The executable will warn
about this file not existing if you don't use it, but that is completely fine.

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
docker compose up --build
# or "docker-compose" on old installs of docker
```

**The executable will warn you if `.env` does not exist**, so when you first run the 
container, please look for that error message and make sure it doesn't occur. If it doesn't,
you can safely pass the `-d` flag to the command

## Testing

This project has unit/integration tests configured. To execute them, you can either run them
locally, or through docker. 

### Locally

```
go test ./...
```

### Through Docker

```
docker compose build .
docker run --rm -it misinfodetector/backend:latest go test ./...
```

See `go help test` for more details for how to run these tests.
