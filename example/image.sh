rm -rf vendor

go get github.com/YAOHAO9/pine

go mod vendor

docker build . -t pine

docker tag pine 192.168.43.126:5000/pine

docker push 192.168.43.126:5000/pine