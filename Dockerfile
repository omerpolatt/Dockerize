# Boyut küçültme işlemi ile beraber Multi-stage Build işlemini gerçekleştirmiş oluyoruz 
#( With the size reduction process, we realize the Multi-stage Build process. )

# alpine ile images boyut küçültme işlemi yapılmıştır ( images size reduction with alpine )
FROM golang:1.22.3-alpine AS builder 


WORKDIR /app

COPY go.mod ./
RUN go mod download


COPY . .


RUN go build -o /app/main .

FROM alpine:latest

#http istekleri için gerekli olan kısım 
RUN apk --no-cache add ca-certificates

WORKDIR /app

COPY --from=builder /app/main .
COPY --from=builder /app/templates ./templates

EXPOSE 8081

CMD ["./main"]
