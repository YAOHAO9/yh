FROM golang

COPY ./ /opt/pine/

WORKDIR /opt/pine/

CMD go run ${server_kind}/main.go 

###################
# Run 'go mod vendor' first 
###################