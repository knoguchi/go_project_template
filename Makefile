VERSION=$(shell git describe --tags)
OBJECT=myproject

default:: all

$(OBJECT):
	go build -ldflags "-X version.GitCommit=$(VERSION)" -o $(OBJECT) cmd/mycmd/main.go

all: $(OBJECT)
