package event

import (
	"encoding/binary"
	"bytes"
	"bufio"
	//"github.com/hraban/opus"
	"io/ioutil"
)

type TransportInfo struct {
	Boundary string
	Message bytes.Buffer
}

func NewTransportInfo( boundary string ) *TransportInfo{
	info := &TransportInfo{}
	info.Boundary = boundary
	return info
}

func ( t * TransportInfo ) Reset() {
	t.Message.Reset()
}

func ( t * TransportInfo ) getJsonHeader( audio bool ) string {
	var content string

	content += "--" + t.Boundary + "\n"

	if audio { // FLAG_JSON_AUDIO_HEADER
		content += "Content-Disposition: "  + "form-data; name=\"audio\"\n"
		content += "Content-Type: " + "application/octet-stream\n"
		content += "\r\n"
	} else {
		content += "Content-Disposition: " + "form-data; name=\"metadata\"\n"
		content += "Content-Type: " + "application/json; charset=UTF-8\n"
		content += "\r\n"
	}

	return content
}

func ( t * TransportInfo ) getJsonBody( event []byte , finished bool ) []byte {
	var content bytes.Buffer
	content.Write( event )
	content.Write( []byte("\n") )

	if finished { //FLAG_JSON_FINISHED
		content.WriteString("--" + t.Boundary + "--")
	}

	return content.Bytes()
}

func ( t * TransportInfo ) CreateMessage( event string ) *TransportInfo{
	t.Reset()
	t.Message.WriteString( t.getJsonHeader( false ) )
	t.Message.Write( t.getJsonBody( []byte( event ) , true ) )
	return t
}

func ( t * TransportInfo ) CreateMessageWithAudioContent( event string , audio []byte ) *TransportInfo{
	t.Reset()
	t.Message.WriteString( t.getJsonHeader( false ) )
	t.Message.Write( t.getJsonBody( []byte( event ) , false ) )
	t.Message.WriteString( t.getJsonHeader( true ) )
	t.Message.Write( t.getJsonBody( audio , true ) )
	return t
}

func ( t * TransportInfo ) CreateAudio( pcm []int16 ) []byte {
	var buffer bytes.Buffer
	binary.Write( bufio.NewWriter(&buffer) , binary.LittleEndian,  pcm )
	return buffer.Bytes()
}

/*
	Decode Opus bytes to PCM bytes
 */
func DecodeOpusToPCM( opusBytes []byte ) ( pcm []int16 , samples int ) {
	//var frameSizeMs float32 = 60
	//dec , _ := opus.NewDecoder(16000, 1 )
	//frameSize := 1 * frameSizeMs * 16000 / 1000
	//pcm = make([]int16, int( frameSize ) )
	//samples, _ = dec.Decode( opusBytes , pcm )
	return pcm , samples
}

/*

 */
func ReadPCMFile( file string )( pcm []int16 ) {
	fileBytes , _ := ioutil.ReadFile( file )
	pcm = make( []int16 , len( fileBytes ) / 2 )
	binary.Read( bytes.NewReader( fileBytes ) , binary.LittleEndian , pcm )
	return pcm
}





