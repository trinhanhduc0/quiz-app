# Sử dụng hình ảnh Go chính thức
FROM golang:1.23 as builder

# Đặt thư mục làm việc
WORKDIR /app

# Sao chép go.mod và go.sum trước
COPY go.mod go.sum ./

# Tải các phụ thuộc
RUN go mod download

# Sao chép mã nguồn
COPY . .

# Biên dịch ứng dụng
RUN go build -tags netgo -ldflags '-s -w' -o cmd/api/main ./cmd/api

# Bước cuối cùng, sử dụng hình ảnh nhỏ hơn để chạy ứng dụng
FROM gcr.io/distroless/base
COPY --from=builder /app/cmd/api/main /main
ENTRYPOINT ["/main"]
