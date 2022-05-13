OBJS = anatool
TEST = *.csv TestResult Include Function *.json

all: anatool

anatool: main.go
	go build -o anatool main.go

clean:
	rm -rf $(OBJS)

cleanTest:
	rm -rf $(TEST)

cleanAll:
	rm -rf $(OBJS) $(TEST)
