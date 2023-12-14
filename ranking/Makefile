BINARY := server
DB_PATH := /tmp/rankings.db

build: compile init_db

compile:
	go build  -o ./$(BINARY) ./main.go

init_db:
ifeq (,$(wildcard $(DB_PATH)))	
	sqlite3 $(DB_PATH) < ./internal/schema.sql
endif

run:
	./$(BINARY)

clean:
	rm -f ./$(BINARY)

