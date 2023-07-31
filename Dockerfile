FROM golang:1.20

RUN git clone https://github.com/blablatov/syslogpars.git
WORKDIR syslogpars

RUN go mod download

COPY *.go ./
COPY *.conf ./

RUN CGO_ENABLED=0 GOOS=linux go build -o /syslogpars

EXPOSE 51444/udp

CMD ["/syslogpars"]