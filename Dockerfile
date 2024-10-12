FROM ubuntu:20.04 as builder                                                                                                                                                                                             
MAINTAINER Maxim Styushin

RUN apt-get update &&\ 
    apt-get install -qq -y wget make ca-certificates git &&\ 
    wget -q -O - https://go.dev/dl/go1.22.5.linux-amd64.tar.gz | tar -C /usr/local -xzf - && ln -s /usr/local/go/bin/go /usr/local/bin/go

WORKDIR /build
COPY go.sum go.mod Makefile .git ./
COPY cmd cmd
COPY pkg pkg
RUN make build

FROM ubuntu:20.04
RUN apt-get update &&\ 
    apt-get install -qq -y ca-certificates &&\
    apt-get clean all &&\
    rm -rf /var/lib/apt/lists/* /var/cache/apt/*

WORKDIR /app
COPY --from=builder /build/bin/go-news-api-gw ./

EXPOSE 8080

CMD ["/app/go-news-api-gw"]
