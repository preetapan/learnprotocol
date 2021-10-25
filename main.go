package main

import (
	"encoding/binary"
	"fmt"
	"net"
	"os"
)

func main() {
	// Initialize operator functions

	var opFns []func(uint64, uint64) uint64
	opFns = append(opFns, add)
	opFns = append(opFns, subtract)
	opFns = append(opFns, multiply)
	opFns = append(opFns, algebra)

	// Listen for incoming connections.
	l, err := net.Listen("tcp", "localhost:7070")
	if err != nil {
		fmt.Println("Error listening:", err.Error())
		os.Exit(1)
	}
	// Close the listener when the application closes.
	defer l.Close()
	fmt.Println("Listening.. ")
	for {
		// Listen for an incoming connection.
		conn, err := l.Accept()
		if err != nil {
			fmt.Println("Error accepting: ", err.Error())
			os.Exit(1)
		}
		// Handle connections in a new goroutine.
		go handleRequest(conn, opFns)
	}
}

// Handles incoming requests.
func handleRequest(conn net.Conn, opFns []func(uint64, uint64) uint64) {
	// Make a buffer to hold incoming data.
	buf := make([]byte, 12)
	// Read the incoming connection into the buffer.
	reqLen, err := conn.Read(buf)
	if err != nil {
		fmt.Println("Error reading:", err.Error())
	}
	// first byte in the buffer is the operator
	operator := buf[0]
	// read length of the two varint encoded operands
	len1 := buf[1]
	len2 := buf[2]

	fmt.Println("Length of two operands ", len1, len2)
	data1 := buf[3 : 3+len1]
	data2 := buf[3+len1 : reqLen]
	// Varint decode the two operands
	operand1, _ := binary.Uvarint(data1)
	operand2, _ := binary.Uvarint(data2)

	// Calculate the result
	answer := opFns[operator](operand1, operand2)

	// Send a response back with the result, varint encode the answer
	resultBuffer := make([]byte, 8)
	len := binary.PutUvarint(resultBuffer, answer)
	fmt.Printf("result encoded in %v bytes\n", len)

	// Write the answer to the connection
	conn.Write(resultBuffer[0:len])
	// Close the connection
	conn.Close()
}

func add(a uint64, b uint64) uint64 {
	return a + b
}

func subtract(a uint64, b uint64) uint64 {
	return a - b
}

func multiply(a uint64, b uint64) uint64 {
	return a * b
}

func algebra(a uint64, b uint64) uint64 {
	return 2*a + b
}
