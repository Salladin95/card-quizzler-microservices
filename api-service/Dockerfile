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
#COPY --from=build /go/src/app/main apiApp
#COPY .env ./
#EXPOSE 8080
#CMD ["./apiApp"]

FROM golang:alpine
RUN mkdir /app
COPY apiApp /app
COPY .env ./
CMD [ "/app/apiApp"]
