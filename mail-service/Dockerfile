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
#COPY --from=build /go/src/app/main mailApp
#COPY .env ./
#CMD ["./mailApp"]

FROM golang:alpine
RUN mkdir /app
COPY mailApp /app
COPY .env ./
CMD [ "/app/mailApp"]
