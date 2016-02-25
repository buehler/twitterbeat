FROM golang:1.5.2
MAINTAINER Christoph Buehler <christoph.buehler@bluewin.ch>

COPY . /go/src/github.com/buehler/go-elastic-twitterbeat

RUN cd /go/src/github.com/buehler/go-elastic-twitterbeat \
    go get -t -u -f && \
    go test -v ./... && \
    go build -o twitterbeat

RUN mkdir -p /etc/twitterbeat/ /var/twitterbeat/data /var/twitterbeat/config && \
    cp /go/src/github.com/buehler/go-elastic-twitterbeat/twitterbeat /etc/twitterbeat/ && \
    cp /go/src/github.com/buehler/go-elastic-twitterbeat/etc/twitterbeat.yml /var/twitterbeat/config/dockerbeat.yml

VOLUME /var/twitterbeat/data /var/twitterbeat/config

WORKDIR /etc/dockerbeat

ENV PERIOD=60 \
    SCREEN_NAMES="[\"@smartive_ch\", \"@elastic\"]" \
    ES_HOSTS="[\"elasticsearch:9200\"]" \
    CONSUMER_KEY="" \
    CONSUMER_SECRET="" \
    ACCESS_KEY="" \
    ACCESS_SECRET=""

CMD [ "./twitterbeat", "-c", "/var/twitterbeat/twitterbeat.yml", "-p", "/var/twitterbeat/data/twitterMap.json" ]