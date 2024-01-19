# Instrumentation Demo App

This is a demo app that we will use to complete challenges for the software instrumentation course.

## About this project

This project is a gRPC server, written in Go, that provides few endpoints to manage courses catalog and bookings for a fake online course marketplace. This project uses library named [gRPC-gateway](https://github.com/grpc-ecosystem/grpc-gateway) which provide a RESTful API endpoints proxy. So, while the implementation is actually a gRPC server, you can also access it by using `curl` or any other way you used to.

## Getting Started

1. Make sure you have terminated all running containers from other projects if any. Run `docker compose down` on the other project if necessary.

1. From this project directory, start all dependencies with docker compose for this project:

    ```bash
    make bootstrap
    docker compose up -d
    ```

    Please check the `docker-compose.yml` file to see what services are started.

1. Seed the database by running the following command from this project directory

    ```bash
    make course/seed
    ```

    If necessary, you may update the data after it is seeded to database. Or you can truncate the database and re-seed it if necessary.

1. To start this API server, run:

    ```bash
    make course/server
    ```
    
    You should see that HTTP server starts at port 8800 and gRPC server starts at port 9900.

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

1. To start the load generator, go to `http://localhost:8089`, fill the parameter as you need, and click `Start Swarming`.