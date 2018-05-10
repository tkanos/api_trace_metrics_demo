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
make deploy
```

### Run
```bash
cd $GOPATH/src/github.com/Tkanos/api_trace_metrics_demo
UID=`id -u` docker-compose up
```

