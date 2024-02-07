# Build Stage
FROM golang:alpine AS build
WORKDIR /go/src/app
COPY . .
RUN go mod download
RUN go build -o main .
# Final Stage
FROM alpine
WORKDIR /app
COPY --from=build /go/src/app/main userApp
COPY .env ./
CMD ["./userApp"]
