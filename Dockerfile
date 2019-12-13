FROM golang

WORKDIR /workspace

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 go build -o /workspace
RUN echo "nobody:x:65534:65534:Nobody:/:" > /etc_passwd

FROM alpine

COPY --from=0 /workspace /workspace
COPY --from=0 /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=0 /etc_passwd /etc/passwd

USER nobody

CMD [ "/workspace/auto-tweet-liker" ]