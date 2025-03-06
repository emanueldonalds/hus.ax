# Hus.ax

https://hus.ax/

View properties listed on Åland.

Built with Go, templ and HTMX.

An API is also available: https://github.com/emanueldonalds/property-api

Site analytics here: https://analytics.edonalds.com/hus.ax

## Requirements
- Go 1.22.2
- [templ 0.2.524](https://templ.guide/)
- docker

## Running locally
Start the test database

```
docker build . -t property-viewer-mariadb 
docker run -d --name property-viewer-mariadb -p 3306:3306 property-viewer-mariadb
```

Then run with env variables:

```
cd src

PROPERTY_VIEWER_DB_HOST=localhost \
PROPERTY_VIEWER_DB_PASSWORD=abc123 \
templ generate --watch --cmd 'go run .'
 ```
