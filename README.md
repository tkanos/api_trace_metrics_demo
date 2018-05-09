# api_trace_metrics_demo

### Requirements

* docker
* docker-compose
* golang

### Setup
```bash
go get -u github.com/Tkanos/api_trace_metrics_demo
cd $GOPATH/src/github.com/Tkanos/api_trace_metrics_demo/app
make install
make build
make deploy
```

### Run
```bash
cd $GOPATH/src/github.com/Tkanos/api_trace_metrics_demo
docker-compose up
```

