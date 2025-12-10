# HTTP-from-TCP-go

This project is implementation a go implementation of HTTP Server built on raw TCP Sockets, no net/http, no frameworks. This project was realized with [Boot.dev course](https://www.boot.dev/courses/learn-http-protocol-golang).

## Project goals and motivation
* To learn and understand how HTTP works at the TCP / socket level (parsing requests, handling connections, writing responses).
* To understand better the underlying logic of HTTP Protocol, and how to implement that in Go.

## What it does?
* Listens on a TCP socket directly, and accepts raw TCP connections.
* Parses incoming bytes into HTTP request according to [RFC 9110](https://datatracker.ietf.org/doc/html/rfc9110).
* Constructs and sends responses manually.
* Handles concurrency using Go goroutines.
