.PHONY: run stop logs new-job

run: 
	docker compose up --build -d

stop:
	docker compose down

logs:
	docker compose logs -f producer consumer1 consumer2

new-job: 
	go run ./cmd/cli \
	--command "echo hello world $$RANDOM" \
	--schedule "*/5 * * * * *"