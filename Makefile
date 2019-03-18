BINS := ComputeLows

all:
	go build ./...

clean:
	rm -rf ${BINS}
