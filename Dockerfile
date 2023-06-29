# Multi-stage syslogpars build
# Многоэтапная сборка syslogpars

FROM golang AS build

ENV location /go/src/github.com/blablatov/syslogpars

WORKDIR ${location}/syslogpars

ADD ./syslogpars.go ${location}/syslogpars

RUN go mod init github.com/blablatov/syslogpars/syslogpars

RUN CGO_ENABLED=0 go build -o syslogpars

# Go binaries are self-contained executables. Используя директиву FROM scratch - 
# Go образы  не должны содержать ничего, кроме одного двоичного исполняемого файла.

FROM scratch
COPY --from=build ./syslogpars ./syslogpars

ENTRYPOINT ["./syslogpars"]
EXPOSE 50051