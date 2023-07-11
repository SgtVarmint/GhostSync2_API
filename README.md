# Requirements

- Golang:
    This project is currently built on 1.20.5

    [Download page](https://go.dev/dl/)

    [Installation Instructions](https://go.dev/doc/install)

# Building the project
## Filling out the env file
Before you can run or build the project, you are going to need to add values to an `env.yml`.

You can either copy-paste the `env.yml.template` file and have both, or rename the `...template` file and fill in the values. It doesn't really matter either way.

Just make sure that `env.yml` is in the project root.

<br>

## Building and running from an exe file
For Linux, in the project root, run the command `go build -C main ../gs2_api` and an binary will appear that you can run.

For Windows, you would need to make sure the file you're writing is an `.exe` file.

<br>

## Running the project without an exe file
In the project root, run the command `go run main/main.go` and the server should start up.