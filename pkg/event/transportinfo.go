package event

import (
	"encoding/binary"
	"bytes"
	"bufio"
	"io"
	"fmt"
	"github.com/hajimehoshi/go-mp3"
	"github.com/hraban/opus"
	"io/ioutil"
)

type MP3Reader struct {
	io.ReadCloser
	buffer bytes.Buffer
}

func ( m * MP3Reader) Read(p []byte) (n int, err error){
	return m.buffer.Read( p )
}

func ( m * MP3Reader) Close() error {
	return nil
}

func NewMP3Reader( data []byte ) * MP3Reader {
	mp3reader := &MP3Reader{}
	mp3reader.buffer.Write( data )
	return mp3reader
}

type PCMWriter struct {
	io.Writer
	io.Reader
	Buffer bytes.Buffer
}

func ( m * PCMWriter ) Read(p []byte) (n int, err error){
	return m.Buffer.Read( p )
}

func ( w * PCMWriter ) Write(p []byte) (n int, err error) {
	return w.Buffer.Write( p )
}

func NewPCMWriter() * PCMWriter {
	return &PCMWriter{}
}

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
	var frameSizeMs float32 = 60
	dec , _ := opus.NewDecoder(16000, 1 )
	frameSize := 1 * frameSizeMs * 16000 / 1000
	pcm = make([]int16, int( frameSize ) )
	samples, _ = dec.Decode( opusBytes , pcm )
	return pcm , samples
}

func EncodePCMToOpus( pcm []int16 ) []byte {
	const bufferSize = 1000 // choose any buffer size you like. 1k is plenty.
	const channels = 1
	enc , _ := opus.NewEncoder( 48000 , 1 , opus.AppVoIP )
	data := make([]byte, bufferSize)
	n, err := enc.Encode(pcm, data)
	if err != nil {
		fmt.Println( err.Error() )
	}

	data = data[:n] // only the first N bytes are opus data. Just like io.Reader.
	return data
}

/*
	Decode MP3 to PCM
	Encode PCM to Opus

	input: mp3
	output: opus
*/
func EncodeMP3ToOpus( mp3data []byte ) []byte {
	var pcm []int16
	dec , _ := mp3.NewDecoder( NewMP3Reader( mp3data ) )
	w := NewPCMWriter()
	io.Copy( w , dec )
	binary.Read( w , binary.LittleEndian , pcm )
	return EncodePCMToOpus( pcm )
}


/*

 */
func ReadPCMFile( file string )( pcm []int16 ) {
	fileBytes , _ := ioutil.ReadFile( file )
	pcm = make( []int16 , len( fileBytes ) / 2 )
	binary.Read( bytes.NewReader( fileBytes ) , binary.LittleEndian , pcm )
	return pcm
}





