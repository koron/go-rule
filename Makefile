test:
	go test -v ./...

tags:
	gotags -f tags -R .

.PHONY: test tags
