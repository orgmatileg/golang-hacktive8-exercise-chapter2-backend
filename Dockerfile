# Luqmanul Hakim / arlhba@gmail.com

# Step 1 membuat binary
FROM golang:alpine AS builder

# Install git karena go get membutuhkan git fetch sebagai depedensi
RUN apk update && apk add --no-cache git

# Mengganti working directory (kalau di linux/mac seperti command cd)
WORKDIR $GOPATH/src/myapp/

# Melakukan copy file dari folder saat ini ke folder working directory
COPY . .

# eksekusi go get untuk mendapatkan semua library luar yang kita gunakan
RUN go get -d -v

# Melakukan build binary apps 
RUN go build -o /go/bin/exercise

# Step 2 - membuat image baru hanya untuk running apps kita dari hasil build di atas
# ini begunakan agar image container kita size nya kecil
FROM alpine

# Melakukan copy binary dari hasil build image sebelumnya ke image scratch ini
COPY --from=builder /go/bin/exercise /exercise

# Melakukan eksekusi binary apps. goodluck!
ENTRYPOINT ["/exercise"]
