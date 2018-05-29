package directive

import (
	"encoding/json"
)

type DirectiveHeader struct {
	Namespace string `json:"namespace"`
	Name string `json:"name"`
	MessageId string `json:"messageId"`
	DialogReqId string `json:"dialogRequestId"`
}

type DirectivePayload struct {
	Url string `json:"url"`
	Format string `json:"format"`
	Token string `json:"token"`
}

type Directive struct {
	Header DirectiveHeader `json:"header"`
	Payload DirectivePayload `json:"payload"`
	mp3data []byte
}

func ( d * Directive ) HasMP3Data() bool {
	if len( d.mp3data ) != 0 {
		 return true
	}
	return false
}

func ( d * Directive ) GetMP3Data() []byte {
	return d.mp3data
}

const (
	PARTTYPE_APPLICATION_JSON = 0
	PARTTYPE_APPLICATION_OCTET_STREAM = 1
)

type Parts struct {
	header []byte
	body []byte
	partType int
}

type DirectiveReader struct {
	response []byte // complete avs response
	boundary string // multipart boundary from avs
	parts []Parts
	jsonParts int
	octetStreamParts int
}

func NewDirectiveReader( raw []byte , boundary string ) []Directive {
	reader := &DirectiveReader{}
	reader.response = raw
	reader.boundary = boundary
	return reader.ReadAll()
}

func ( r * DirectiveReader ) analyseContent() {

}

func ( r * DirectiveReader ) ReadAll() []Directive {

	r.analyseContent()

	directives := make( []Directive , r.jsonParts )
	directiveCount := 0

	for _ , jsonPart := range r.parts {
		if jsonPart.partType == PARTTYPE_APPLICATION_JSON {
			json.Unmarshal( jsonPart.body , &directives[ directiveCount ] )
		}
	}

	return directives
}
