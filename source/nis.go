package main

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"fmt"
	"net"
	"time"
)

type NetworkInformationServer struct {
	Connection net.Conn
}

func ( networkInformationServer *NetworkInformationServer ) Connect( address net.IP, port int, timeout int ) ( err error ) {
	connection, connectError := net.DialTimeout( "tcp4", fmt.Sprintf( "%s:%d", address, port ), time.Duration( timeout ) * time.Millisecond )
	if connectError != nil { return connectError }

	networkInformationServer.Connection = connection

	return nil
}

func ( networkInformationServer *NetworkInformationServer ) Disconnect() ( err error ) {
	disconnectError := networkInformationServer.Connection.Close()
	if disconnectError != nil { return disconnectError }

	return nil
}

func ( networkInformationServer *NetworkInformationServer ) SendCommand( command string ) ( bytesSent int, err error ) {
	var buffer bytes.Buffer

	// Length as 16-bit big-endian unsigned integer
	writeLengthError := binary.Write( &buffer, binary.BigEndian, uint16( len( command ) ) )
	if writeLengthError != nil { return 0, writeLengthError }

	// Command as raw bytes
	_, writeCommandError := buffer.Write( []byte( command ) )
	if writeCommandError != nil { return 0, writeCommandError }

	// Send
	bytesSent, sendError := networkInformationServer.Connection.Write( buffer.Bytes() )
	if sendError != nil { return 0, sendError }

	return bytesSent, nil
}

func ( networkInformationServer *NetworkInformationServer ) ReceiveResponse() ( response string, err error ) {
	connectionReader := bufio.NewReader( networkInformationServer.Connection )
	var buffer bytes.Buffer

	// TODO: Look at first 24 bytes for information about response (e.g. 'APC : 001,036,0857')

	for {

		// Length
		lengthBytes := make( []byte, 2 )
		readLengthError := binary.Read( connectionReader, binary.BigEndian, lengthBytes )
		if readLengthError != nil { return "", readLengthError }
		dataLength := binary.BigEndian.Uint16( lengthBytes )

		// EOL
		if dataLength == 0 { break }

		// Data
		dataBytes := make( []byte, binary.BigEndian.Uint16( lengthBytes ) )
		readDataError := binary.Read( connectionReader, binary.BigEndian, dataBytes )
		if readDataError != nil { return "", readDataError }

		// Append
		_, appendError := buffer.Write( dataBytes )
		if appendError != nil { return "", appendError }

	}

	return buffer.String(), nil
}

func ( networkInformationServer *NetworkInformationServer ) FetchStatus() ( status Status, err error ) {
	// Send
	_, sendError := networkInformationServer.SendCommand( "status" )
	if sendError != nil { return Status{}, sendError }

	// Receive
	statusResponse, receiveError := networkInformationServer.ReceiveResponse()
	if receiveError != nil { return Status{}, receiveError }

	// Parse
	status, parseError := ParseStatusText( statusResponse )
	if parseError != nil { return Status{}, parseError }

	return status, nil
}
