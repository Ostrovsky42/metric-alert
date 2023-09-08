	AGENT_VERSION=v1.20
	SERVER_VERSION=v1.20
ci:
	golangci-lint run -v -c golangci.yaml

test:
	go test -cover ./...

analyze:
	go run cmd/staticlint/main.go -source ./...

build: 	build/agent build/server

build/agent:
	go build -o ./cmd/agent/agent \
	-ldflags "-X main.buildVersion=$(AGENT_VERSION) \
		  -X main.buildDate=$(shell date '+%H:%M:%S[%Y/%m/%d]') \
		  -X 'main.buildCommit=$(shell git log --pretty=format:"%h  %s" -n 1)'" \
	./cmd/agent/*.go

build/server:
	go build -o ./cmd/server/server \
	-ldflags "-X main.buildVersion=$(SERVER_VERSION) \
		  -X main.buildDate=$(shell date '+%H:%M:%S[%Y/%m/%d]') \
		  -X 'main.buildCommit=$(shell git rev-parse HEAD)'" \
	./cmd/server/*.go

profile/agent:
	go tool pprof -http=":9091" -seconds=30 http://localhost:6061/debug/pprof/profile

profile/server:
	go tool pprof -http=":9090" -seconds=30 http://localhost:6060/debug/pprof/profile

doc:
	godoc -http=:8090 -play

swagger:
	swag init --output ./swagger/ -g internal/server/handlers/update.go -g internal/server/handlers/get_value.go

auto: build
	metricstest -test.v -test.run=^TestIteration14$ \
                                      -agent-binary-path=cmd/agent/agent \
                                      -binary-path=cmd/server/server \
                                      -file-storage-path=/tmp/metrics-db.json \
                                      -server-port=8080 \
                                      -key="jojo" \
                                      -database-dsn='postgres://metric_user:metric_pass@localhost:5433/metric_alert?sslmode=disable' \
                                      -source-path=.
