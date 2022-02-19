# fuck-off-as-a-service


## Overview

This repository runs a server that gets messages from [fuck off as a service](https://www.foaas.com/) and returns them. The
server implements a [sliding window log rate limiter](https://bhargav-journal.blogspot.com/2020/12/understanding-rate-limiting-algorithms.html).

## Requirements

The server is implemented in [Go 1.15](https://go.dev). 
To install Go, follow the [instructions](https://go.dev/doc/install).

## Getting Started

- To run the server, go to the root folder and execute:

```
make all
make run
```

`make all` command will test and build the server.

`make run` command will start the server.

If you have docker installed and want to run it in a containerised environment, run:

```
make docker-build
make docker-run
```

After setting up the server, it can be tested with the following `curl`. 
The request must contain a header called `User-Id`

Example:

```
curl localhost:8080/message -H 'User-Id: "Santiago Leira"'
```

- To test the code and see the coverage, go to the root folder and execute:

```
make test
make coverage
```

## Customising the server

There are many arguments to customize the server:

- log-level, by default it's debug. Ex: it can be info. 
- rate-limit-count, by default it's 5. It's the number of requests allowed in the window time.
- rate-limit-window-in-milliseconds, by default it's 10000. It's the window time to evaluate the number of requests. 
- timeout-in-milliseconds, by default it's 10000. It's the timeout of the API call to `fuck off as a service`.

Example:

```
./fuck-off-as-a-service serve --log-level=info rate-limit-count=1 rate-limit-window-in-milliseconds=20000 timeout-in-milliseconds=20000
```

## Future work

- Implement authorization.
- Support different algorithms as rate limiter.
- Use a database (ex: Redis) to support the rate limiter in a distributed system.  
- Improve the tests (add more acceptance and integration tests).
- Improve errors (add more custom errors, be more precise with errors).
- Improve the documentation on the code.
- Improve CI/CD.
