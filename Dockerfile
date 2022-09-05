FROM golang:1.17.0-alpine3.14 AS builder

LABEL stage=gobuilder

ENV CGO_ENABLED 0
ENV GOOS linux
ENV GOARCH amd64

WORKDIR /build

RUN apk --no-cache add vips vips-dev tzdata

COPY . .

RUN go mod download

RUN go mod verify

RUN go build -ldflags="-w -s" -o golang-api

FROM alpine:3.14

COPY --from=builder /usr/share/zoneinfo/Africa/Cairo /usr/share/zoneinfo/Africa/Cairo
ENV TZ Africa/Cairo

WORKDIR /app

RUN apk --no-cache add vips vips-dev 

RUN apk --no-cache add tini

ENV APP_ENV production

ENV UID=10001

RUN addgroup -S golang-api-service

RUN adduser -D \    
	--disabled-password \    
	--gecos "" \    
	--home "/nonexistent" \    
	--shell "/sbin/nologin" \    
	--no-create-home \    
	--uid "${UID}" \    
	golang-api-user \ 
	-G golang-api-service

USER golang-api-user

COPY --chown=golang-api-user:golang-api-service --from=builder /build/golang-api /app/golang-api

ENTRYPOINT [ "/sbin/tini", "--" ]

CMD ["/app/golang-api"]