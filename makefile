
all: gogl

gogl:
	go build

run: clean gogl
	./GoGL

clean:
	$(RM) GoGL
