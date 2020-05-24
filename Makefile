.PHONY: snapshot dist test vet lint fmt clean
OUT := republik-feeder
PKG := github.com/maetthu/republik-feeder
PKG_LIST := $(shell go list ${PKG}/... | grep -v /vendor/)
GO_FILES := $(shell find . -name '*.go' | grep -v /vendor/)

all: snapshot

snapshot:
	goreleaser --snapshot --skip-publish --rm-dist

dist:
	goreleaser --rm-dist

test:
	@go test -v ${PKG_LIST}

vet:
	@go vet ${PKG_LIST}

lint:
	@for file in ${GO_FILES} ;  do \
		golint $$file ; \
	done

fmt:
	@gofmt -l -w -s ${GO_FILES}

clean:
	-@rm ${OUT}
