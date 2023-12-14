FROM docker:latest

RUN apk add --no-cache go make sqlite python3-dev py-pip 
RUN apk add --no-cache sdl2-dev dpkg-dev freetype-dev #Environment requirements
COPY games/requirements.txt .
RUN python3 -m pip install -r requirements.txt
COPY . /sources
WORKDIR /sources
ENV GOCACHE=/root/.cache/go-build
RUN --mount=type=cache,target="/root/.cache/go-build" make

ENTRYPOINT make run
