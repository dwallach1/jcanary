FROM golang:1.18-alpine

WORKDIR /app
ENV PATH /go/bin:$PATH

COPY go.mod ./
COPY go.sum ./
RUN go mod download
COPY . .
RUN go build -o /jcanary

# to do -- create slimmed down version

ENTRYPOINT ["/jcanary"]