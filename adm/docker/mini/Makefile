all: test

imgname=todogomini

build:
	docker build -f Dockerfile -t ${imgname} .

test:
	docker run --rm -it ${imgname} make -C todogo test

demo:
	docker run --rm -it ${imgname} todo config -i

clean:
	rm -f *~

clean-img:
	docker image rm ${imgname}

clean-all: clean clean-img
