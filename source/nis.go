package main

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"fmt"
	"net"
	"time"
)

// Structure to hold the TCP connection and functions
type NetworkInformationServer struct {
	Connection net.Conn
}

// Connects to a Network Information Server
func ( networkInformationServer *NetworkInformationServer ) Connect( address net.IP, port int, timeout int ) ( err error ) {
	
	// Try to connect using TCP
	connection, connectError := net.DialTimeout( "tcp4", fmt.Sprintf( "%s:%d", address, port ), time.Duration( timeout ) * time.Millisecond )
	if connectError != nil { return connectError }

	// Update the structure property
	networkInformationServer.Connection = connection

	// Return no errors
	return nil

}

// Disconnects from a Network Information Server
func ( networkInformationServer *NetworkInformationServer ) Disconnect() ( err error ) {

	// Try to close the connection
	disconnectError := networkInformationServer.Connection.Close()
	if disconnectError != nil { return disconnectError }

	// Return no errors
	return nil

}

// Sends a command to the Network Information Server
func ( networkInformationServer *NetworkInformationServer ) SendCommand( command string ) ( bytesSent int, err error ) {
	
	// Create an empty buffer
	var buffer bytes.Buffer

	// Add the command length as 16-bit big-endian unsigned integer
	writeLengthError := binary.Write( &buffer, binary.BigEndian, uint16( len( command ) ) )
	if writeLengthError != nil { return 0, writeLengthError }

	// Add the command as raw bytes
	_, writeCommandError := buffer.Write( []byte( command ) )
	if writeCommandError != nil { return 0, writeCommandError }

	// Send the command to the server
	bytesSent, sendError := networkInformationServer.Connection.Write( buffer.Bytes() )
	if sendError != nil { return 0, sendError }

	// Return the number of bytes sent
	return bytesSent, nil

}

// Receives a full response from the Network Information Server
func ( networkInformationServer *NetworkInformationServer ) ReceiveResponse() ( response string, err error ) {
	
	// Create a reader and an empty buffer
	connectionReader := bufio.NewReader( networkInformationServer.Connection )
	var buffer bytes.Buffer

	// TODO: Look at first 24 bytes for information about response (e.g. 'APC : 001,036,0857')

	// Loops until the end of the response is reached
	for {

		// Parse the length as 16-bit big-endian unsigned integer
		lengthBytes := make( []byte, 2 )
		readLengthError := binary.Read( connectionReader, binary.BigEndian, lengthBytes )
		if readLengthError != nil { return "", readLengthError }
		dataLength := binary.BigEndian.Uint16( lengthBytes )

		// Stop if we reached the end of the response
		if dataLength == 0 { break }

		// Extract the remaining data
		dataBytes := make( []byte, binary.BigEndian.Uint16( lengthBytes ) )
		readDataError := binary.Read( connectionReader, binary.BigEndian, dataBytes )
		if readDataError != nil { return "", readDataError }

		// Add data to the endn of the buffer
		_, appendError := buffer.Write( dataBytes )
		if appendError != nil { return "", appendError }

	}

	// Return the data in the buffer as text
	return buffer.String(), nil

}

// Helper function to send the status command and give the response in a structure
func ( networkInformationServer *NetworkInformationServer ) FetchStatus() ( status Status, err error ) {

	// Send the status command
	_, sendError := networkInformationServer.SendCommand( "status" )
	if sendError != nil { return Status{}, sendError }

	// Receive the response
	statusResponse, receiveError := networkInformationServer.ReceiveResponse()
	if receiveError != nil { return Status{}, receiveError }

	// Parse the response
	status, parseError := ParseStatusText( statusResponse )
	if parseError != nil { return Status{}, parseError }

	// Return the status structure
	return status, nil

}
