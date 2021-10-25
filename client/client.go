package client

import (
	"bufio"
	"encoding/binary"
	"fmt"
	"net"
)

func CalculatorClient() {
	data := []uint64{20, 127, 128, 256, 32768, 131072, 67108866, 1048576}
	operator := []byte{0, 2, 3, 1}
    index := 0
	for i:=0; i < len(data); i+=2 {
		c, err := net.Dial("tcp", "localhost:7070")
		if err != nil {
			fmt.Println(err)
			return
		}
		msgBuffer := make([]byte, 3)

		// First write the operator to the buffer
		msgBuffer[0] = operator[index]
		index+=1

		// Encode the two operands
		data1 := encodeVarInt(data[i])
		data2 := encodeVarInt(data[i+1])

		// Write length of operands to the message buffer
		msgBuffer[1] = byte(len(data1))
		msgBuffer[2] = byte(len(data2))

		// Write varint encoded operands to the message buffer
		msgBuffer = append(msgBuffer, data1...)
		msgBuffer = append(msgBuffer, data2...)

		fmt.Println("Data sent to server..")
		for _, x := range msgBuffer {
			fmt.Printf("%v ", x)
		}
		fmt.Println()
		// Send the message to the server and read response back
		c.Write(msgBuffer)

		// Read results from the server
		resultBuffer := make([]byte, 8)
		len, err := bufio.NewReader(c).Read(resultBuffer)
		if err != nil {
			fmt.Println("ERROR! ", err)
			continue
		}

		// Decode the result with varint
		answer, errata := binary.Uvarint(resultBuffer[0:len])
		if errata  < 0 {
			fmt.Printf("Result did not get decoded correctly %v\n", errata)
		}
		fmt.Printf("Result from server: %v\n", answer)
		fmt.Println("************************")
	}

}

func encodeVarInt(data uint64) ([]byte) {
	dataBuffer := make([]byte, 8)
	len := binary.PutUvarint(dataBuffer, data)
	fmt.Printf("%v encoded to length %v\n", data, len)
	// compress slice down to its varint encoded length
	dataBuffer = dataBuffer[0:len]
	return dataBuffer
}
