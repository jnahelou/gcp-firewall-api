FROM golang:1.13 as builder
ENV HOME /app
ENV CGO_ENABLED 0
ENV GOOS linux
WORKDIR /app
ADD . /app/
RUN go get -v -d
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -ldflags '-extldflags "-static"' -o main .

FROM scratch
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /app/main /app/
WORKDIR /app
CMD ["./main"]
