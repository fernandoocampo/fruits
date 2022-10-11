[![Go Report Card](https://goreportcard.com/badge/github.com/fernandoocampo/fruits)](https://goreportcard.com/report/github.com/fernandoocampo/fruits) ![CI](https://github.com/fernandoocampo/fruits/actions/workflows/quality.yaml/badge.svg?branch=master) ![Docker Build](https://github.com/fernandoocampo/fruits/actions/workflows/build.yml/badge.svg?branch=master) ![linter](https://github.com/fernandoocampo/fruits/actions/workflows/static-analysis.yml/badge.svg?branch=master) ![GitHub release](https://img.shields.io/github/v/release/fernandoocampo/fruits.svg?include_prereleases&style=plastic)

# fruits

Service to handle fruits information.

## How to build?

from project folder executes commands below, it will generate binary files in the `./bin/` folder

* just build with current operating system
```sh
make build
```

* build for a linux distro operating system
```sh
make build-linux
```

## How to run a test environment quickly?

1. make sure you have `docker-compose` installed.
2. run the docker compose.
```sh
docker-compose up -d
```

or run this shortcut

```sh
make run-local
```

or

```sh
make test
```

3. once you finished to use the environment, follow these steps

```sh
make clean-local
```

## How to test?

from project folder run the following command

```sh
go test -race ./...
```

## How to call the fruit API

You can use insomnia api client and use the project `insomnia-fruits-service.json`.

## Coding Decisions

1. The service was built following the hexagonal architecture pattern in order to improve maintainability and extensibility. Most of the logic of the service is related to external resources like loggers, databases and monitoring platforms.
2. Loose coupling between packages is very important to increase cohesion, so I avoided referencing another package directly.
3. All service methods must receive the context parameter. The idea is to propagate the cancellation of context and other values in the future. e.g. correlation id.
4. An in-memory database was built to persist fruits information, this was done thinking that this was an mvp and we don't want to use a well-known database engine yet.
5. The adapters folder has all the packages that provide logic to communicate with external components.
6. To keep a loose coupling between the `service` package and the` memorydb` package, I created the `repository` package which provides the logic that both packages need to communicate with each other.
7. `monitoring` package provides logic to monitor and handle metrics for the fruits service and the `metrics` package simulates the integration with an instrumenting platorm where agents will send the metrics.
8. A worker was created to define the `monitoring` logic. The idea is to process metrics asynchronously to avoid adding latency to the core capabilities of the service.
9. I think that the web logic of the service is part of the adapters in the hexagonal architecture pattern, that's why the package `web` is within the `adapter` folder.
10. I used the `go-kit` library to build the microservice in order to make things easy to understand rather than easy to do.
11. In my opinion the main function needs to be very clean, that's why it just calls the application logic to run the service.
12. The `application` package is in charge of start the application, instantiate components and make dependency injections.
13. The `configurations` package provides the logic to load the application setup.
14. The fruit dataset will be loaded when the application starts, in case only one record is invalid, its status will be invalid.
15. Because the status of the dataset is part of the logic of the fruits, the fruit package will take care of indicating whether the dataset is valid or not.
16. Assumption: The mandatory fields for fruit are: name, classification, country and vault.
17. docker-compose was used to build and run the application, so far this is only one service. The service can be run in a standalone fashion though.
18. Please notice that the service has a lot of logs, the idea behind this is to facilitate debugging at production or qa environments.

## Improvements for a live production system.

1. We should provide an endpoint to verify that the service is running or not. i.e `/health` or `/heartbeat`.
2. Just in case security is an important issue for this service, we should provide a security mechanism to allow only granted users to make changes in the fruits service data. e.g. `JWT`
3. We should change the current in-memory database for a well known engine. e.g postgresql or mongodb. Data must survive a service interruption.
4. The monitoring platform should be changed by a well known engine such as Prometheus.
5. It's important to create a pipeline to run unit and integration tests, build and deploy the application in different environments.
6. Generate documentation for the fruit service API in order to facilitate its use. For this we can use Swagger.


## deploy in kubernetes

* only use this with minikube
eval $(minikube docker-env)
docker build -t local/fruits:latest .

https://www.baeldung.com/ops/kubernetes-helm

```sh
helm install --name fruits ./k8s-v2/fruits
```


## using flux

check project [flux-practices](https://github.com/fernandoocampo/flux-practices)