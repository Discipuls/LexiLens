FROM golang:1.22.5

WORKDIR /go/src/app

COPY go.mod go.sum ./

RUN go mod download && go mod verify


COPY . .
RUN go build -v

CMD [ "./LexiLensTgBot" ]


