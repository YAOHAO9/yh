FROM golang

COPY ./* /opt/program/

WORKDIR /opt/program/

RUN go install

CMD go run main.go