# Build and Run

```sh
go build .
./scheduler
```

# Test API

## Add tasks

```sh
curl -X POST http://localhost:8080/add_tasks \
-H "Content-Type: application/json" \
-d '{"tasks": [{"cost": 2}, {"cost": 4}, {"cost": 100}, {"cost": 6}]}'
```

## Set policy

The default policy is FIFO.

set policy to SRTF:

```sh
curl -X POST http://localhost:8080/set_policy \
-H "Content-Type: application/json" \
-d '{"policy": "srtf"}'
```

set policy to FIFO:

```sh
curl -X POST http://localhost:8080/set_policy \
-H "Content-Type: application/json" \
-d '{"policy": "fifo"}'
```

# Run go tests for policy

```sh
go test -run TestFIFOScheduling
go test -run TestSRTFScheduling
```
