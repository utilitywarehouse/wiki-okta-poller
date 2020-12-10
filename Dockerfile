FROM: golang:alpine

WORKDIR /wiki-okta-poller

COPY . .

RUN go get ./...

CMD ["/wiki-okta-poller/bin/graph

