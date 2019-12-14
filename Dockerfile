FROM golang

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .


RUN CGO_ENABLED=0 go build -o /dist/auto-tweet-liker
RUN echo "nobody:x:65534:65534:Nobody:/:" > /etc_passwd

FROM scratch

WORKDIR /app

COPY --from=0 /dist/auto-tweet-liker ./
COPY --from=0 /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=0 /etc_passwd /etc/passwd

USER nobody

ENTRYPOINT ["/app/auto-tweet-liker"]