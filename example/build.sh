docker build -t pine .
docker run --rm --name pine -e server_port=5678 -dt pine