date = $(shell bash ../scripts/getDate.sh)
version = $(shell bash ../scripts/getVersion.sh)

exe = prim2tab
all : $(exe)

$(exe): 
	go build -ldflags "-X github.com/evolbioinf/prim/util.version=$(version) -X github.com/evolbioinf/prim/util.date=$(date)" $(exe).go
$(exe).go: $(exe).org
	awk -f ../scripts/preTangle.awk $(exe).org | bash ../scripts/org2nw | notangle -R$(exe).go | gofmt > $(exe).go
test: $(exe)_test.go $(exe)
	@go test
tangle: $(exe).go $(exe)_test.go
$(exe)_test.go: $(exe).go
	awk -f ../scripts/preTangle.awk $(exe).org | bash ../scripts/org2nw | notangle -R$(exe)_test.go | gofmt > $(exe)_test.go	
clean:
	rm -f *.go $(exe)
