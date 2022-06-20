default: bin/gen-enum
	go run ./input/main.go > ./input/defs.json
	bin/gen-enum --input-file ./input/defs.json
	$(MAKE) clean
.PHONY: default

bin/gen-enum:
	mkdir -p bin
	go build -o $@ ./output

.PHONY: clean
clean:
	rm -rf bin
