FROM golang:1.17.2-alpine AS build
RUN apk add --no-cache git
## CGO_ENABLED=0 要不然编译文件，无法执行 https://blog.csdn.net/u013235478/article/details/105852353/
RUN go env -w GOPROXY=https://goproxy.cn,direct CGO_ENABLED=0
RUN go get github.com/golang/glog

WORKDIR /go/src/project
#RUN dep ensure -vendor-only

COPY . /go/src/project/
RUN go build -o /bin/project

FROM scratch
COPY --from=build /bin/project /bin/project 
ENTRYPOINT ["/bin/project"]
