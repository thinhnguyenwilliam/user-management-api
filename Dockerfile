# ---------- Stage 1: Build ----------
FROM golang:1.25-alpine AS builder

# install git (cần cho go mod)
RUN apk add --no-cache git

WORKDIR /app

# copy go mod trước (optimize cache)
COPY go.mod go.sum ./
RUN go mod download

# copy source code
COPY . .

# build binary (static)
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o user-management-api cmd/api/main.go

# ---------- Stage 2: Run ----------
FROM alpine:3.19

WORKDIR /app

# install tzdata (optional nhưng nên có)
RUN apk add --no-cache tzdata

# copy binary từ builder
COPY --from=builder /app/user-management-api .

# expose port
EXPOSE 8086

# run app
CMD ["./user-management-api"]