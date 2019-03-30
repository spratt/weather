BINS := PlotAverageTemperature

all:
	go build ./...

clean:
	rm -rf ${BINS}
