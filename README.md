# es-redis-proxy

A caching proxy for Elasticsearch requests

## Building

### Mac/Linux

0) set GOROOT environment variable
1) Install Go and Make
2) make

### Docker

0) set GOROOT environment variable
1) Install Docker, Go and Make
2) make docker


## Running

### Mac/Linux

```
./es-redis-proxy
```

### Docker

```
docker pull maguec/es-redis-proxy:latest
docker run -i -t -p 8080:8080 maguec/es-redis-proxy
```

## Testing

run either the docker container or the raw application binary

```
curl http://localhost:8080/health

#check the proxy
curl --header 'Content-Type: application/json' -X POST http://localhost:8080/instruments/_search -d '{"query": {"match_all": {}}}'  |jq
#check the source
curl --header 'Content-Type: application/json' -X POST http://localhost:9200/instruments/_search -d '{"query": {"match_all": {}}}'  |jq

```

---
Copyright © 2021, Chris Mague
