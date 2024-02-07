# Build Stage
FROM golang:alpine AS build
WORKDIR /go/src/app
COPY . .
RUN go mod download
RUN go build -o main .
# Final Stage
FROM alpine
WORKDIR /app
COPY --from=build /go/src/app/main apiApp
COPY .env ./
EXPOSE 8080
CMD ["./apiApp"]
