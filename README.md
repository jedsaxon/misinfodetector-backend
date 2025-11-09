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
    	where this program should listen for api requests (default :memory:)
  -rabbitmq string
    	rabbitmq connection string (default amqp://guest:guest@localhost:5672/)
  -sqlite string
    	where the sqlite database should be stored (default 127.0.0.1:5000)
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

The example contains commented-out environment variables, which you may override by
uncommenting them as needed. If they are left commented or discluded entirely from
this file, then the defaults will be loaded.

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

## Misinformation States

| ID  | Name        |
| --- | ----------- |
| 0   | Fake        |
| 1   | True        |
| 2   | Not Checked |

## Endpoints

### Get Posts

`GET /api/posts`

Retrieves a paginated list of posts from the database. Returns the posts along with the total number of pages available.

**Query Parameters**

- `pageNumber` (integer, required): The page number to retrieve.
    - Must be greater than 0
- `resultAmount` (integer, required): The amount of results to return.
    - Must be greater than 0
    - Must be less than 51

**Example Usage**

```sh
curl -X GET "http://localhost:5000/api/posts?pageNumber=1&resultAmount=10" \
-H "Content-Type: application/json"
```

**Response - HTTP 200**

Will return a paged list of products. Along with all posts in the current page, it will also return a `pages` param
containing the total amount of pages the user can navigate through.

`misinfo_state` is the ID of the misinformation state, which can be found in the table [above](#misinformation-states)

```json
{
  "message": "10 posts found",
  "posts": [
    {
      "id": "uuid-string",
      "message": "Post content here",
      "username": "john_doe",
      "submitted_date": "2025-11-05T10:30:00Z",
      "misinfo_state": {
        "state": 0,
        "confidence": 0.5,
        "submitted_date": "2006-01-02T15:04:05Z"
      }
    }
    }
  ],
  "pages": 5
}
```

**Response - HTTP 404**

If the post could not be found, the server will respond with 404, with the following body

```json
{
  "message": "request contains errors",
  "errors": {
    "id": "Post with the given ID could not be found"
  }
}
```

### Create Post

`POST /api/posts`

Creates a new post with the provided message and username. Once the post is created and stored in the database, it will
be published into the RabbitMQ queue, for misinformation detection processing. The misinformation state by default will
be set to 2, until the misinformation detector service can handle it.

**Request Body**

```json
{
  "message": "This is my post content",
  "username": "john_doe"
}
```

- `username` (string, required): username of the poster.
    - length must be above 0
    - length must be bellow 64
- `message` (string, required): content of the post
    - length must be above 0
    - length must be bellow 256

**Example Usage**

```sh
curl -X POST "http://localhost:5000/api/posts" \
-H "Content-Type: application/json" \
-d '{ "message": "This is my post content", "username": "john_doe" }'
```

**Response - HTTP 201**

If the post was successfully created, the server will respond with 201 created. It will set the `Location` header to the
URL of the created post.

`misinfo_state` is the ID of the misinformation state, which can be found in the table [above](#misinformation-states)

```json
{
  "message": "successfully created post",
  "post": {
    "id": "uuid-string",
    "message": "This is my post content",
    "username": "john_doe",
    "submitted_date": "2025-11-05T10:30:00Z",
    "misinfo_state": {
      "state": 0,
      "confidence": 0.5,
      "submitted_date": "2006-01-02T15:04:05Z"
    }
  }
}
}
```

### Import Posts

`PUT /api/posts`

Imports posts from a CSV file, into the database. It expects the following columns in a csv formatted file (can have different
column names, but must be in correct order):

`id` - Id of the post (not used, but required)
`text` - The post's actual content
`date` - When the post was submitted
`label` -
`pred_label` - Whether the post contains misinformation. `1` = true, `0` = fake
`pred_prob` - Probability that this post contains misinformation. A float between 0-1
`correct` - Whether this record should be included in the import. True/Falset .

**Request Body**

You will need to upload a file with the name `posts`.

**Example Usage**

```sh
curl localhost:5000/api/posts -X PUT \
-F "posts=@/.../path/to/predictions_detailed.csv"
```

**Response - HTTP 204**

If all posts were inserted into the database, it will return HTTP 204 no content.

**Response - HTTP 400**

If you did not upload a file, then it will return HTTP 400

### Put Random Posts

`PUT /api/posts/random`

Creates a specified number of randomly generated posts for testing purposes.

**Request Body**

```json
{
  "amount": 10
}
```

- `amount` (integer, required): number of posts to insert
    - Must be more than 0
    - Cannot be greater than 20,000

**Example Usage**

```sh
curl -X POST "http://localhost:5000/api/posts/random" \
-H "Content-Type: application/json" \
-d '{ "amount": 10 }'
```

### Get TNSE Embeddings

`GET /api/data/tnse-embeddings`

Retrieves all t-SNE (t-Distributed Stochastic Neighbor Embedding) visualization data from the database. This endpoint
returns embedding coordinates and prediction information for visualizing the misinformation detection model's results.

**Query Parameters**

None Required.

**Example Usage**

```sh
sh curl -X GET "http://localhost:5000/api/data/tnse-embeddings" \
-H "Content-Type: application/json"
```

**Response - HTTP 200**

Returns an array of t-SNE embedding records. Each record contains the original label, predicted label, correctness,
and 2D coordinates for visualization.

```json
[
  {
    "id": 0,
    "label": 1,
    "pred_label": 0,
    "correct": "True",
    "tnse_x": 0.00000,
    "tnse_y": 0.00000
  }
]
```

### Import T-SNE Embeddings

Imports t-SNE embedding data from a CSV file into the database. This endpoint will delete all existing embeddings before
importing the new data. The CSV file should contain the following columns in order:

- `id` - Record ID
- `label` - Actual label (0 = fake, 1 = true)
- `pred_label` - Predicted label (0 = fake, 1 = true)
- `correct` - Whether the prediction was correct ("True"/"False")
- `tnse_x` - X coordinate in t-SNE space
- `tnse_y` - Y coordinate in t-SNE space

**Request Body**

You will need to upload a file with the form field name `embeddings`.

**Example Usage**

```sh
sh curl -X PUT "http://localhost:5000/api/data/tnse-embeddings" \
-F "embeddings=@/path/to/tnse_embeddings.csv"
```

**Response - HTTP 204**

If all embeddings were successfully imported into the database, it will return HTTP 204 no content.

**Response - HTTP 400**

If the file was not provided or the form field name is incorrect, the server will respond with 400.

### Get Topic Activities

`GET /api/data/topic activities`

**Query Parameters**

None Required.

**Example Usage**

```sh
sh curl -X GET "http://localhost:5000/api/data/topic-activities" \
-H "Content-Type: application/json"
```

**Response - HTTP 200**

Returns an array of t-SNE embedding records. Each record contains the original label, predicted label, correctness,
and 2D coordinates for visualization.

```json
[
  {
    "db_id": 0,
    "date": "20216-10-27T00:00:00Z",
    "text": 0,
    "topic_id": 0,
    "topic_name": "abcde"
  }
]
```

### Import Topic Activities

CSV formatted document for topic activities.

- `date` should be in format YYYY/MM/DD
- `text` Text containing activity
- `topic_id` The id of the topic
- `topic_name` Name of the topic

**Request Body**

You will need to upload a file with the form field name ``.

**Example Usage**

```sh
sh curl -X PUT "http://localhost:5000/api/data/topic-activities" \ 
-F "topics=@/path/to/topic_activities.csv"
```

**Response - HTTP 204**

If all topics were successfully imported into the database, it will return HTTP 204 no content.

**Response - HTTP 400**

If the file was not provided or the form field name is incorrect, the server will respond with 400.

## Common Error Responses

### 400 Bad Request

Returned when validation fails or request parameters are invalid. It will return all errors for all fields, in a
`Record<string, string>` under the `errors` property.

```json
{
  "message": "request contains errors",
  "errors": {
    "field_name": "field error message"
  }
}
```

### 500 Internal Server Error

Returned when a server-side error occurs. These errors may be out of your control.

```json
{
  "message": "internal server error"
}
```
