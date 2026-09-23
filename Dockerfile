FROM golang:1.23-alpine AS build
WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 go build -o /splitmate ./cmd/splitmate
FROM scratch
COPY --from=build /splitmate /splitmate
EXPOSE 8080
ENTRYPOINT ["/splitmate"]
