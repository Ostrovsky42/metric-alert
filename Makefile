
ci:
	golangci-lint run -v -c golangci.yaml

test:
	go test -cover ./...

build: 	build/agent build/server

build/agent:
	go build -o ./cmd/agent/agent ./cmd/agent/*.go

build/server:
	go build -o ./cmd/server/server  ./cmd/server/*.go

profile/agent:
	go tool pprof -http=":9091" -seconds=30 http://localhost:6061/debug/pprof/profile

profile/server:
	go tool pprof -http=":9090" -seconds=30 http://localhost:6060/debug/pprof/profile



auto: build
	metricstest -test.v -test.run=^TestIteration14$ \
                                      -agent-binary-path=cmd/agent/agent \
                                      -binary-path=cmd/server/server \
                                      -file-storage-path=/tmp/metrics-db.json \
                                      -server-port=8080 \
                                      -key="jojo" \
                                      -database-dsn='postgres://metric_user:metric_pass@localhost:5433/metric_alert?sslmode=disable' \
                                      -source-path=.
