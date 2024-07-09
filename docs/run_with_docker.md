# Using Makefile to Run Verifier and Prover with Docker

Make sure you have Docker, Docker Compose and Make installed on your machine. If not, follow their official guides for
installation. Once done, follow the steps below:

## Building the Docker Images

To build the Docker images for the `verifier` and `prover`, navigate to the project's root directory and execute the
build task from the Makefile:

```bash
make build
```

This will remove any previous images named `n-zkp-test-verifier` and `n-zkp-test-prover` then build new images
for `n-zkp-test-verifier` and `n-zkp-test-prover` using the provided Dockerfiles.

## Running the Docker Containers

You can use the `make up` command to run your applications:

```bash
make up
```

This command will stop any previously running containers and then use Docker Compose to create and start the `verifier`
and `prover` services. It will also follow the logs of the containers.

## Stopping the Docker Containers

If you need to stop the running services you can use the `make down` command:

```bash
make down
```

This command calls Docker Compose to stop the running services defined in your `docker-compose.yml`.

## Configure environment variables

You can modify the initial setup by editing the docker-compose.yml file.