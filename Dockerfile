FROM golang:1.24 AS build

WORKDIR /app
COPY src/go.mod src/go.sum ./
RUN go mod download

COPY src /app
RUN go tool templ generate
RUN CGO_ENABLED=0 GOOS=linux go build -o /entrypoint

FROM gcr.io/distroless/static-debian11 AS final

WORKDIR /
COPY --from=build /entrypoint /entrypoint
COPY --from=build /app/assets /assets
COPY --from=build /app/rss/template.xml /rss/template.xml

EXPOSE 4932
USER nonroot:nonroot
ENTRYPOINT ["/entrypoint"]
