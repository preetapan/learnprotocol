Overview
========
Implementation of a calculator service built on a custom protocol 
on top of TCP 

Details
========
The server is in `main.go`, and the client is in `client/client.go`
Run `go test` from the client package to see a set of end to end examples

Protocol Details
================
The first byte is used to represent an operation code.
The next two bytes represent the length in bytes of two varint encoded operands.
The remaining bytes contain varint encoded bytes for the two operands

More info
=========
https://developers.google.com/protocol-buffers/docs/encoding
https://carlmastrangelo.com/blog/lets-make-a-varint 