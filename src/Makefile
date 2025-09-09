LINTER_PATH=/tmp/bin
LINTER_BIN=/tmp/bin/golangci-lint

install-go:
	wget -c https://go.dev/dl/go1.20.14.linux-amd64.tar.gz \
	&& sudo tar -xvf go1.20.14.linux-amd64.tar.gz -C /usr/local \
	&& rm go1.20.14.linux-amd64.tar.gz

go-version:
	go version

install-lint-deps:
	(which golangci-lint > /dev/null) || curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(LINTER_PATH) v1.55.2

lint-version: install-lint-deps
	$(LINTER_BIN) version

lint-fix: install-lint-deps
	$(LINTER_BIN) run --fix ./...

lint: install-lint-deps
	$(LINTER_BIN) run ./...

.PHONY: install-go go-version lint-version lint-fix lint
