FROM golang:latest

ENV HOME=/home/server
ENV GO_PORT=8085

WORKDIR ${HOME}

COPY go.mod ${HOME}/
COPY go.sum ${HOME}/
RUN go mod download
RUN go install github.com/cosmtrek/air@latest

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o /docker-gs-ping

EXPOSE ${GO_PORT}

CMD air