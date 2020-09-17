# meant4

## What is it?

Test task for meant4.com

Compiled and checked with go 1.13.6.

## Building, testing and running

Build: `go build .`  
Test: `go test .`  
Run locally: `go run .`  
Dockerize: `docker build . --tag bhmj`  
Run in docker: `docker run -p 8989:8989 bhmj`  

## Details

The resulting numbers are strings since it is a safe way to represent such large values.