.PHONY: default
default: test megacheck vet golint gocyclo coverage

.PHONY: test
test:
	go test -v -race .

.PHONY: coverage
coverage:
	go test -race -coverprofile=coverage.txt -covermode=atomic

.PHONY: megacheck
megacheck:
	megacheck .

.PHONY: vet
vet:
	go vet .

.PHONY: golint
golint:
	golint -set_exit_status .

.PHONY: gocyclo
gocyclo:
	gocyclo -over 10 .

.PHONY: install
install:
	go install .

.PHONY: deps
deps:
	go get github.com/golang/lint/golint
	go get honnef.co/go/tools/cmd/megacheck
	go get github.com/fzipp/gocyclo