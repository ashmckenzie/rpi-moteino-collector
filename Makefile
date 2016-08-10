RELEASES_DIR=$(realpath -e ${PWD})/releases

build:
	docker build -t moteino-collector:latest .

releases: build
	docker run --rm -v ${RELEASES_DIR}:/out moteino-collector:latest rsync -ax bin/ /out/

run: build
	docker run --rm -ti moteino-collector:latest

shell: build
	docker run --rm -ti moteino-collector:latest bash
