# Project Structure

This document explains the structure of your project's directory. It implements Domain-Driven Design principles.

## Directory Structure

1. **`cmd`**: Hosts the main applications for the project, it currently has two applications named `prover`
   and `verifier`. Each contains a `main.go` file, serving as the entry point for the respective applications.
2. **`config`**: Holds files related to the project configuration. This could include loading and parsing configuration
   files.
3. **`Docker-related files` (`docker-compose.yml`, `Dockerfile.prover`, and `Dockerfile.verifier`)**: Files for creating
   Docker images and orchestrating application containers. `docker-compose.yml` defines the multi-container application
   setup, `Dockerfile.prover` and `Dockerfile.verifier` define the Docker images creation steps for `prover`
   and `verifier` respectively.
4. **`integration_tests`**: Stores integration test which verify the interactions both components.
5. **`internal`**: Designed to store packages for use within this project exclusively. Overview of its directories:
    - **`app`**: Holds the core business logic of the application. The magic happens here.
    - **`domain`**: Holds domain entities. `auth` handles authentication-related logic.
    - **`interactor`**: Manages interactivity between other layers, like transforming data from the repository layer for
      presentation layer use.
    - **`repository`**: Data access layer responsible for interaction with the persistence layer (database, in-memory
      data store etc).
6. **`proto`**: Holds Protocol Buffer files, used for serializing structured data for data exchange across
   different services or components.

