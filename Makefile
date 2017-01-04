# dependencies that are used by the build & test process
DEPEND=github.com/Masterminds/glide

all: clean depend build

clean:
	rm -f video_server/run_server && rm -f game_server/run_server

build: video game

video:
	go build -o video_server/run_server video_server/*.go

game:
	go build -o game_server/run_server game_server/*.go

# installing build dependencies. You will need to run this once manually when you clone the repo
depend:
	go get -u -v $(DEPEND)
	${GOPATH}/bin/glide install
