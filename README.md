# Instrumentation Demo App

This is a demo app that we will use to complete challenges for the software instrumentation course.

## About this project

This project is a gRPC server, written in Go, that provides few endpoints to manage courses catalog and bookings for a fake online course marketplace. This project uses library named [gRPC-gateway](https://github.com/grpc-ecosystem/grpc-gateway) which provide a RESTful API endpoints proxy. So, while the implementation is actually a gRPC server, you can also access it by using `curl` or any other way you used to.

## Getting Started

### Pre-setup

I have prepared some basic setup here including the following:

* Docker compose file to start few dependencies.
* Logger initialization which will write the logs to `stdout` and file located on `logs/app.log`. Thus, you can just directly use zerolog to log your application. You can find the initialization on [server.go](./cmd/course/commands/server.go#L51).

### Running the application

1. Make sure you have terminated all running containers from other projects if any. Run `docker compose down` on the other project if necessary.

1. Set these environment variable in your profile (e.g. `~/.bashrc`):

    ```bash
    export REDIS_PASSWORD=<set to password you like>
    ```

1. From this project directory, start all dependencies with docker compose for this project:

    ```bash
    make bootstrap
    docker compose up -d
    ```

    Make sure you have `make` installed. If not, run `apt install make` or whatever you need to do to install `make` on your machine.
    
    If you get error when running `make bootstrap` because some dependencies are missing, just install them. For instance, if you got this following errors:

    ```
    no required module provides package google.golang.org/grpc/cmd/protoc-gen-go-grpc; to add it:
        go get google.golang.org/grpc/cmd/protoc-gen-go-grpc
    ```

    In this case, run `go get google.golang.org/grpc/cmd/protoc-gen-go-grpc`.

    Please check the `docker-compose.yml` file to see what services are started.

1. To start this API server, run:

    ```bash
    make course/server
    ```
    
    You should see that HTTP server starts at port 8800 and gRPC server starts at port 9900.

1. Seed the database by running the following command from this project directory

    ```bash
    make course/seed
    ```

    If necessary, you may update the data after it is seeded to database. Or you can truncate the database and re-seed it if necessary.

1. Check out list of available APIs from the swagger docs. Go to `http://localhost:8800/swagger`. You can try out the API from there as well if you want.

## Running load generator

Load generator here is used to test your implementation. 

While you may see the implementation of the load generator, I suggest you to **not do that**. Instead, you should treat the load generator as a black box so that you can fully understand how the API clients interacting with your API **only by using instrumentation** you set.

1. Make sure you have python 3.x installed. 

1. Install virtualenv

    ```bash
    pip install virtualenv
    ```

1. Go to `locust` directory under `scripts` directory. Crete and activate a virtual environment

    ```bash
    cd scripts/locust
    virtualenv .venv
    source .venv/bin/activate
    ```

1. Install the required dependencies for the load generator

    ```bash
    pip install -r requirements.txt
    ```

1. Run the load generator

    ```bash
    locust -f locustfiles --class-picker --modern-ui -H http://localhost:8800
    ```    

1. Open your browser and go to `http://localhost:8089`. You should see the locust dashboard. There are two users scenarios:

    * `GeneralUser` is the scenario that you can use to test several APIs. Use this for most of the time when working on the final challenges. You do not need to spawn a lot of users to test this scenario. One or just a few users are enough.

    * `CompetingUser` is the scenario that you can use later when working on the tracing final challenge. The scenario is designed to simulate a competition between users when making the reservation. You need to spawn a more users to test this scenario to see how it might impact the performance of the system.

1. To start the load generator, go to `http://localhost:8089`, fill the parameter as you need, and click `Start Swarming`.