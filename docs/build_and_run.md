# Building and Running the Verifier and Prover

Follow the steps below to build and run the `verifier` and `prover` applications. These applications reside within
the `cmd` directory of your project.

## Prerequisites

Ensure that Go (version 1.22 or later) is installed on your machine, and `GOPATH` environment variable is configured
correctly.

## Instructions

### **Building the Verifier and Prover Applications**

Navigate to the root directory of your project and run these commands to build the `verifier` and `prover`:

```bash
go build -o verifier ./cmd/verifier
go build -o prover ./cmd/prover
```

These commands will create executable files `verifier` and `prover` in your project's root directory.

### **Running the Verifier and Prover Applications**

Now, you can run the `verifier` and `prover` applications with (do it in 2 different consoles and make sure you spin up
the verifier):

```bash
./verifier
./prover
```