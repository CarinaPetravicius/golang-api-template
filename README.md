# golang-api-template
API Template with Golang, Hexagonal Architecture, OpenApi, Prometheus, Postgres, Redis, Kafka and OpenPolice.

### Install Go
`sudo apt-get update`  
`sudo apt-get -y upgrade`  
`wget  https://go.dev/dl/go1.21.0.linux-amd64.tar.gz`  
`sudo tar -xvf go1.21.0.linux-amd64.tar.gz`  
`sudo mv go /usr/local`  
`export GOROOT=/usr/local/go`  
`export GOPATH=$HOME/Projects/golang-api-template`  
`export PATH=$GOPATH/bin:$GOROOT/bin:$PATH`  
`go version`

### Install Docker and docker-compose
- To install Docker follow the official documentation: https://docs.docker.com/engine/install/ubuntu/
- To install docker compose: `sudo apt install docker-compose`
- Give permission to execute: `sudo chmod 666 /var/run/docker.sock`

### How to install an updated dependency
- This example is a dependency to read yaml files: 
`go get gopkg.in/yaml.v3`

### Start the project
- In the root of this project, run to start the database: docker-compose up
- Install all the dependencies defined on go.mod file: `go get .`
- Config the environment variables of your localhost environment: `DATABASE_DNS=postgresql://root:root@localhost:5432/db?sslmode=disable`

### Kafka interface
After run the project, you can access the Kafdrop on:
- `localhost:9000` 

On this Kafka interface you can see that the kafka topic was created.

### Endpoints for Health Check:
- `http://localhost:8080/health/live`
- `http://localhost:8080/health/ready`

### Prometheus endpoint with Go and Http metrics with custom service_name label:
- `http://localhost:8080/metrics`
