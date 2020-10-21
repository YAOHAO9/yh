FROM golang

COPY ./ /opt/pine/

WORKDIR /opt/pine/

# RUN go mod download ###### go mod vendor

CMD go run ${Server_Kind}/main.go 