# Tracking Server Demonstration

Demonstration project for tracking users impressions when browsing through the web.

## Overview

This project creates a web server application which captures, parses and stores HTTP requests' body contents into [Apache Cassandra](https://cassandra.apache.org/) databases. The web server is built using the [Go](http://golang.org/) programing language and shipped as a [Docker](https://www.docker.com/) container to [Kubernetes](https://kubernetes.io/) clusters - optimizing cost and performance.

The web server responds to the following endpoints:

| Endpoint | Method |              |
|----------|--------|--------------|
| /        | GET    | Healtcheck   |
| /track   | POST   | Data storing |

For the `/track` endpoint, the application will expect a json containing the internet event data using the following pattern:

```json
{
    "username": "John Doe",
    "target": "https://www.google.com/",
    "description": "VIEW"
}
```

There are two valid values for the `description` field: _CLICK_ and _VIEW_. Any other values will trigger a `BadRequest` response, informing that the field's value is invalid.

## Deployment

Both the application and database of this project are shipped as containers to Kubernetes clusters. One can recreate this infrastructure by using the `.yaml` files present in the `k8s` directory. The deploy creates a `StorageClass`, a `PersistentVolumeClaim`, and a `StatefulSet` with a headless `Service` for the database as well as a `Deployment` and `LoadBalancer` for the API.

There is a demonstration version of this project running in a cluster provided by [DigitalOcean](https://www.digitalocean.com/) in the address `178.128.130.123:8080`. To test it you can run the following commands:

`curl 178.128.130.123:8080/`

`curl -X POST 178.128.130.123:8080/track -d '{"username": "John Doe", "target": "https://www.google.com/", "description": "VIEW"}'`

### Handling High Throughput of Tracking Events

Since we are dealing with the "whole" web, we can expect a high throughput of requests in our API. In the provided backend, there is a load balancer that can help us with that. In case the pods behind the load balancer aren't able to handle the incoming traffic, we can easily scale our application by increasing the `replicas` number of our Cassandra `StatefulSet` and of our `Deployment`.

## Local Development

You can build and run the application locally using [docker-compose](https://docs.docker.com/compose/). Install it and run the following command on the project's root directory:

`docker-compose up --build`

Even though the output is kind of strange, by the end there should be an local app listening to the `8080` port of your machine and a Cassandra standalone cluster listening to the port `9042` of your machine. At this point, if you want, you can run the unit and integration tests present in the project:

`go test -v ./...`


## Possible Improvements
* Enable TLS connections
* Logging and monitoring
