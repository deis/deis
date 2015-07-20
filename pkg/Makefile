include ../includes.mk

repo_path = github.com/deis/deis/pkg

GO_PACKAGES = prettyprint time
GO_PACKAGES_REPO_PATH = $(addprefix $(repo_path)/,$(GO_PACKAGES))

test: test-style test-unit

test-style:
# display output, then check
	$(GOFMT) $(GO_PACKAGES)
	@$(GOFMT) $(GO_PACKAGES) | read; if [ $$? == 0 ]; then echo "gofmt check failed."; exit 1; fi
	$(GOVET) $(GO_PACKAGES_REPO_PATH)
	$(GOLINT) ./...

test-unit:
	$(GOTEST) ./...
