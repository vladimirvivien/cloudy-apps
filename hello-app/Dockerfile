FROM golang:1.14.0 as gobuilder
WORKDIR /go-app
COPY . /go-app
ENV GO111MODULE=on
RUN CGO_ENABLED=0 go build -a -tags netgo -ldflags '-w -extldflags "-static"' -o hello-app

FROM scratch
COPY --from=gobuilder /go-app/hello-app .
ENTRYPOINT ["/hello-app"]