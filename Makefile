BINARY=bin/ufocatch
SOURCEDIR=.
SOURCES := $(shell find $(SOURCEDIR) -name '*.go')
TESTS := $(shell go list ./... | grep -v 'vendor')

.DEFAULT_GOAL:$(BINARY)

$(BINARY): $(SOURCES)
	go build -o ${BINARY} .

.PHONY: clean test install

install:
	go install .

clean:
	if [ -f ${BINARY} ]; then rm ${BINARY}; fi

test:
	go test -v -race $(TESTS)
