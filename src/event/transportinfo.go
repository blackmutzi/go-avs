package event

type TransportInfo struct {
	Boundary string
	Message string
}

func NewTransportInfo( boundary string ) *TransportInfo{
	info := &TransportInfo{}
	info.Boundary = boundary
	info.Message = ""
	return info
}

func ( t * TransportInfo ) Reset() {
	t.Message = ""
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

func ( t * TransportInfo ) getJsonBody( event string , finished bool ) string {
	var content string
	content += event + "\n"

	if finished { //FLAG_JSON_FINISHED
		content += "--" + t.Boundary + "--"
	}

	return content
}

func ( t * TransportInfo ) CreateMessage( event string ) *TransportInfo{
	t.Reset()
	t.Message += t.getJsonHeader( false )
	t.Message += t.getJsonBody( event , true )
	return t
}

func ( t * TransportInfo ) CreateMessageWithAudioContent( event string , audio []byte ) *TransportInfo{
	t.Reset()
	t.Message += t.getJsonHeader( false )
	t.Message += t.getJsonBody( event , false )
	t.Message += t.getJsonHeader( true )
	t.Message += t.getJsonBody( string( audio ) , true )
	return t
}




