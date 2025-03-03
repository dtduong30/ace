# Stage 1: Build
# Check version trong go.mod
FROM golang:1.23-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o wishlist cmd/main.go

# Stage 2: Run
FROM alpine:latest
WORKDIR /app
# Chỉ copy file cần chạy, bỏ lại Go compiler và source code, giữ image nhẹ
COPY --from=builder /app/wishlist .
# Nếu đổi port khác phải sửa ở r.Run trong main
EXPOSE 8080
# Dùng [] để chạy trực tiếp, không qua shell (tối ưu hơn)
CMD ["./wishlist"]

# Note: Nếu chỉ dùng FROM golang:1.21-alpine, image cuối sẽ nặng (~300MB), chứa cả compiler, source code không cần thiết

