# Hus.ax

https://hus.ax/

View properties listed on Åland.

Built with Go, templ and HTMX.

Site analytics here: https://analytics.edonalds.com/hus.ax

## Requirements
- Go >= 1.24
- docker

## Running locally
Start the test database:

```
docker compose up db-local -d
```

Install dependencies:

```
go mod download
```

Then run with env variables:

```
cd src

PROPERTY_VIEWER_DB_HOST=localhost \
PROPERTY_VIEWER_DB_PASSWORD=abc123 \
go tool templ generate --watch --cmd 'go run .'
 ```

To run everything in docker:

```
docker compose --profile local up -d
```
