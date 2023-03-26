## Build
FROM golang:1.19 AS build

ENV GOPROXY=https://mirrors.cloud.tencent.com/go/

WORKDIR /app

# Download necessary Go modules
COPY [ "go.mod", "go.sum", "./" ]
RUN go mod download

# Copy and build the source code
COPY . .
RUN go build -o /simplified-tik-tok

## Deploy
FROM gcr.io/distroless/base-debian11

WORKDIR /

COPY --from=build /simplified-tik-tok /simplified-tik-tok

EXPOSE 8888

USER nonroot:nonroot

CMD [ "/simplified-tik-tok" ]
