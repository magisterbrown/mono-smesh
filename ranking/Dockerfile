FROM golang:latest

COPY . /sources
WORKDIR /sources
ENV GOCACHE=/root/.cache/go-build
RUN mv ./config/paths_deploy.go ./config/paths.go
RUN --mount=type=cache,target="/root/.cache/go-build" go build  -o ./server ./main.go

CMD ["./server"]
