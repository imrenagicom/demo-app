# Load Generator

## Running the load generator

1. Install virtualenv

    ```bash
    pip install virtualenv
    ```

1. Create a virtual environment

    ```bash
    virtualenv .venv
    ```

1. Activate the virtual environment

    ```bash
    source .venv/bin/activate
    ```

1. Install the dependencies

    ```bash
    pip install -r requirements.txt
    ```

1. Run the load generator

    ```bash
    locust -f locustfiles --class-picker --modern-ui -H http://localhost:8800
    ```    


