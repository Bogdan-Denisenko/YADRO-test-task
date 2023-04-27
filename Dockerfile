FROM golang:1.19-alpine3.17

WORKDIR /app

RUN mkdir -p ./testData
RUN mkdir -p ./cmd
RUN mkdir -p ./internal

COPY go.mod ./
COPY internal ./internal
COPY cmd ./cmd
COPY testData ./testData

RUN go build -o ComputerClub ./cmd/main.go

CMD ./ComputerClub $FILE_PATH