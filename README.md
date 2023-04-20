# Shorten App
Shorten App is a simple web application that allows you to shorten URLs using a hash function and store them in a Redis database. It also provides an endpoint to retrieve the original URL from the shortened URL.

## The programming language used in this application is Go (also known as Golang)
Some of the pros of using Go include:

- __High performance and low memory footprint__: Go's design focuses on efficiency, making it a great choice for building high-performance applications, especially those with large-scale concurrent operations.
- __Built-in concurrency__: Go has a built-in concurrency model called Goroutines, which allows for lightweight parallelism and efficient use of system resources.
- __Cross-platform support__: Go can be compiled and run on multiple operating systems, making it a versatile choice for building applications that need to run on different environments.
- __Strong type safety__: Go's type system is designed to catch errors at compile-time, reducing the likelihood of runtime errors.

## Requirements
Docker
Docker Compose
Postman or similar HTTP client

## Usage
Clone the repository:

```
git clone https://github.com/your-username/shorten-app.git
```

Navigate to the project directory:

```
cd shorten-app
```

Start the containers:

```
docker-compose up
```

Open Postman or a similar HTTP client and test the endpoints (see below).
When you are done, stop the containers by pressing CTRL+C in the terminal.

## Endpoints
### GET /redis
Returns information about the Redis server, including version, uptime, memory usage, and so on.

Example:

```
curl -X GET http://localhost:8080/redis
```

### GET /health
Returns a simple message indicating that the Shorten App is healthy.

Example:

```
curl -X GET http://localhost:8080
```

### POST /shorten
Shortens a given URL and returns the shortened URL.

Request body:

```
{
    "long_url": "https://www.example.com/some/long/path/to/a/page"
}
```

Response body:

```
{
    "long_url": "https://www.example.com/some/long/path/to/a/page",
    "short_url": "http://localhost:8080/7e3d3df939"
}
```

Example:

```
curl -X POST -H "Content-Type: application/json" -d '{"long_url": "https://www.example.com/some/long/path/to/a/page"}' http://localhost:8080/shorten
```

### GET /:short_url
Redirects the user to the original URL associated with the given shortened URL.

Example:

```
curl -X GET http://localhost:8080/7e3d3df939
```

Notes
The shortened URLs are valid for 24 hours only, after which they are automatically deleted from the Redis database.
The application runs on port 8080 by default. If you need to change the port, you can modify the docker-compose.yml file accordingly.
The Redis server runs on port 6379 by default. If you need to change the Redis configuration, you can modify the docker-compose.yml file accordingly.



