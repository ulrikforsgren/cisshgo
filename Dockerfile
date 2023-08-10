FROM golang:1.17.1-alpine as build
ADD . /app
WORKDIR /app
ENV GO111MODULE=on
RUN go mod download
RUN go build -o /app/cisshgo cissh.go


FROM alpine as nx
RUN mkdir /app
COPY --from=build /app/cisshgo /app/cisshgo
COPY /transcripts /app/transcripts
COPY /id_rsa /app/.
COPY /start.sh /app/.
WORKDIR /app
ENTRYPOINT ["/app/start.sh"]
