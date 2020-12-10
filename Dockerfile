FROM golang:alpine

WORKDIR /wiki-okta-poller

COPY . .

RUN apk update && apk upgrade && apk add --no-cache bash git

RUN export GOPATH=`pwd` && go get ./...

CMD ["/wiki-okta-poller/bin/wiki-poller"]

