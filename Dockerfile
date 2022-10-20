FROM golang:1.18-alpine as BUILD

WORKDIR /src/
COPY . /src/

RUN go mod tidy
RUN CGO_ENABLED=0 go build -o /bin/app .

FROM alpine
COPY --from=BUILD /bin/app /bin/app
ENTRYPOINT [ "/bin/app" ]
