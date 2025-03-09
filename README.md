# Go Url Shortener

-   [go-chi](https://github.com/go-chi/chi) for routing
-   [shortuuid](https://github.com/lithammer/shortuuid) for generating UUIDs

## Usage

```bash
> go run main.go

# On another terminal
> http --form POST :3000/shorten url=https://www.google.com

http://localhost:3000/short/<some-uuid>

# Go to the shortened URL and you will be redirected to the original URL
```

## Load testing

-   [Locust](https://locust.io/) is used for load testing

```bash
uvx locust --config locust.conf
```
