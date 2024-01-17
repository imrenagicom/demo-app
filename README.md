# Instrumentation Demo App

## Getting Started

1. Make sure you have terminated all running containers from other projects if any. Use `docker compose down`.

1. From this project directory, start all dependencies with docker compose for this project:

    ```bash
    docker compose up -d
    ```

    Please check the `docker-compose.yml` file to see what services are started.

1. Seed the database by running the following command from this project directory

    ```bash
    make course/seed
    ```

1. To start this API server, run:

    ```bash
    make course/server
    ```
    
    You should see that HTTP server starts at port 8800 and gRPC server starts at port 9900.

1. To try out the API, go to `http://localhost:8800/swagger` and try out the endpoints.