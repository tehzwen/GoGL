
all: gogl

gogl:
	go build

run:
	./gogl

clean:
	$(RM) gogl