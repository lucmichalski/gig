.PHONY: clean test
build: gig

test:
	go test -race -v ./... -count=1

integ-test:
	go test -race -covermode atomic -coverprofile=profile.cov -v ./... --tags=integration -count=1

lint:
	docker run --rm -v $$(pwd):/app -w /app golangci/golangci-lint:latest golangci-lint run -v
	
gig:
	go build .

clean:
	rm gig
