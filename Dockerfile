FROM golang:1.5.2
MAINTAINER Christoph Buehler <christoph.buehler@bluewin.ch>

ENV GO15VENDOREXPERIMENT=1

RUN git clone https://github.com/Masterminds/glide.git $GOPATH/src/github.com/Masterminds/glide && \
    cd $GOPATH/src/github.com/Masterminds/glide && \
    make bootstrap && \
    make build && \
    cp ./glide /usr/bin

COPY . $GOPATH/src/github.com/buehler/twitterbeat

RUN cd /go/src/github.com/buehler/twitterbeat && \
    make build-docker

RUN mkdir -p /etc/twitterbeat/ /var/twitterbeat/data /var/twitterbeat/config && \
    cp /go/src/github.com/buehler/twitterbeat/twitterbeat /etc/twitterbeat/ && \
    cp /go/src/github.com/buehler/twitterbeat/etc/twitterbeat.yml /var/twitterbeat/config/twitterbeat.yml

RUN rm -rf /go

VOLUME /var/twitterbeat/data /var/twitterbeat/config

WORKDIR /etc/twitterbeat

ENV PERIOD="60s" \
    SCREEN_NAMES="[\"@smartive_ch\", \"@elastic\"]" \
    ES_HOSTS="[\"elasticsearch:9200\"]" \
    CONSUMER_KEY="" \
    CONSUMER_SECRET="" \
    ACCESS_KEY="" \
    ACCESS_SECRET=""

CMD [ "./twitterbeat", "-c", "/var/twitterbeat/config/twitterbeat.yml", "-p", "/var/twitterbeat/data/twitterMap.json", "-e", "-v" ]