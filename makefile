
all: gogl

gogl:
	go build

run: clean gogl
	./gogl

clean:
	$(RM) gogl