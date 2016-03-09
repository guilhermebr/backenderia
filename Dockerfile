FROM golang:1.6-wheezy

RUN go get github.com/tools/godep

WORKDIR $GOPATH/src/github.com/guilhermebr/backenderia

ADD . $GOPATH/src/github.com/guilhermebr/backenderia

RUN godep restore

CMD ["go", "run", "main.go"]
