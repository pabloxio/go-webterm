BIN_DIR = bin
BIN_NAME = term
build: $(BIN_DIR) $(BIN_DIR)/$(BIN_NAME)

$(BIN_DIR):
	@mkdir -p $(BIN_DIR)

GOSOURCES = $(shell find . -type f -name "*.go")
$(BIN_DIR)/$(BIN_NAME): $(GOSOURCES)
	@go build -o $(BIN_DIR)/$(BIN_NAME)

test:
	@go test -cover ./...

GITIGNORE ?= go
gitignore:
	curl -Ls "http://www.gitignore.io/api/$(GITIGNORE)" | tee .gitignore
	@if [ -f .gitignore.custom ]; then \
		cat .gitignore.custom >> .gitignore; \
	fi

clean:
	@rm -rf $(BIN_DIR)/

.PHONY: clean
