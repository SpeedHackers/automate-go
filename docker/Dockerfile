FROM golang

MAINTAINER Josh Chase <jcjoshuachase@gmail.com>
RUN go get github.com/SpeedHackers/automate-go/automated && rm -rf /go/pkg /go/src && strip /go/bin/automated
EXPOSE 8888
EXPOSE 8444

ADD start.sh /start.sh
CMD ["/bin/bash", "/start.sh"]
