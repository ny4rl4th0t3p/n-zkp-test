# Running Tests in The Project

## **Running All Tests**

To run all the tests within the project (except for the integration tests), navigate to your project's root directory
and run:

```bash
go test ./... -v --count=1
```

This command will test all the packages in your project by recursively testing all subdirectories.

## **Running Integration Tests**

The project have separate integration tests in `integration_tests` directory, you could run test here separately using
the `go test` command:

```bash
go test -tags integration ./integration_tests/integration_test.go 
```
