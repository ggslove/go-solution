FROM golang:1.16-alpine AS build
RUN apk add --no-cache git
RUN go env -w GOPROXY=https://goproxy.cn,direct
RUN go get github.com/golang/glog

COPY Gopkg.lock Gopkg.toml /go/src/project/
WORKDIR /go/src/project/
RUN dep ensure -vendor-only

COPY . /go/src/project/
RUN go build -o /bin/project

FROM scratch
COPY --from=build /bin/project /bin/project ENTRYPOINT ["/bin/project"]
CMD ["--help"]