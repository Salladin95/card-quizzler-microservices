## Build Stage
#FROM golang:alpine AS build
#WORKDIR /go/src/app
#COPY . .
#RUN go mod download
## Build the Go application located in the cmd/api directory and name the output executable 'main'
#RUN go build -o main ./cmd/api/
## Final Stage
#FROM alpine
#WORKDIR /app
#COPY --from=build /go/src/app/main userApp
#COPY .env ./
#CMD ["./userApp"]

FROM alpine:latest
RUN mkdir /app
COPY userApp /app
COPY .env ./
CMD [ "/app/userApp"]
