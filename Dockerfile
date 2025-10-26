FROM golang:1.25

WORKDIR /usr/src/app

# pre-copy/cache go.mod for pre-downloading dependencies and only redownloading them in subsequent builds if they change
COPY go.mod go.sum ./
RUN go mod download
RUN mkdir /usr/local/bin -p
RUN mkdir /var/lib/backend -p

COPY . .
RUN go build -v -o /usr/local/bin/ ./... 

CMD ["misinfodetector-backend"]

