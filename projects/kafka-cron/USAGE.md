- Use `docker compose` to start the service.
    ```shell
    docker compose up -d
    ```

- Use cli to add a job.
    ```shell
    go run ./cmd/cli --schedule "0 0 * * *" --command "echo hello"
    ```
