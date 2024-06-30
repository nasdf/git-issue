.PHONY: docs
docs:
	go run -tags docs .

.PHONY: install
install:
	go install .