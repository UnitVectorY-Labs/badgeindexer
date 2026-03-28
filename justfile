
# Commands for badgeindexer
default:
  @just --list
# Build badgeindexer with Go
build:
  go build ./...

# Run tests for badgeindexer with Go
test:
  go clean -testcache
  go test ./...