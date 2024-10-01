FROM golang:1.23-alpine as builder

WORKDIR /app
COPY . /app

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o codenotary-be ./cmd/codenotarybe

FROM alpine

EXPOSE 8000

COPY --from=builder /app/codenotary-be .
CMD [ "./codenotary-be" ]