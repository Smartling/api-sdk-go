.PHONY: tidy
tidy:
	go mod tidy

.PHONY: test_unit
test_unit:
	go test ./api/... ./helpers/...
