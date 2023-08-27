# Overview

This file will have the guidelines to use the application

## Docker

In order to run the application in a Docker container, should first create the image with this command:

> docker build -t servercounter .

This command can be improve with more information

Than, you can create the container with the following command:

> docker run --publish 8080:8080 servercounter

## Unit tests

This projects have unit tests and in order to run it you run the command:

> go run test -v

This command it's executed when you create the docker image to guarantee that everything it's ok before creating

## Load tests

You can test the application with the help of [Vegeta](https://github.com/tsenart/vegeta) tool.

After create the Docker container and running on port 8080, you can run the command:

> echo "GET http://localhost:8080" | vegeta attack -duration=5s -rate=50 | vegeta report

This command will send 50 request in 5 seconds to our application. You can change those values to test
