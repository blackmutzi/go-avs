package directive

import (
	"encoding/json"
	"strings"
	"bytes"
)

const (
	PARTTYPE_APPLICATION_JSON = 0
	PARTTYPE_APPLICATION_OCTET_STREAM = 1
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

type ResponseDirective struct {
	Directive Directive `json:"directive"`
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

type MultiPart struct {
	ContentID string
	ContentType string
	Content []byte
	MultpartType int
}

type DirectiveReader struct {
	response []byte // complete avs response
	boundary string // multipart boundary from avs
	parts []MultiPart
	jsonParts int
	octetStreamParts int
}

func NewDirectiveReader( raw []byte , boundary string ) []ResponseDirective {
	reader := &DirectiveReader{}
	reader.response = raw
	reader.boundary = boundary
	return reader.ReadAll()
}

func ( r * DirectiveReader ) analyseContent() {
	parts := bytes.Split( r.response[ len( []byte( r.boundary ) ) : ] , []byte( r.boundary ) )
	r.parts = make( []MultiPart , len( parts ) - 1 )

	for count , partContent := range parts {
		if strings.HasPrefix( string( partContent ) , "--"){
			break
		}

		if cidpos := bytes.Index(partContent, []byte("Content-ID:") ); cidpos != -1 {

			ctypepos := bytes.Index(partContent,  []byte("Content-Type:") )
			coctetpos := bytes.Index(partContent,  []byte("application/octet-stream") ) + len(string("application/octet-stream"))
			id3pos := bytes.Index(partContent,  []byte("ID3") )

			r.parts[ count ].ContentID = string( partContent[ cidpos: cidpos + ctypepos - 2 ] )
			r.parts[ count ].ContentType = string ( partContent[ ctypepos: ctypepos + coctetpos ] )
			r.parts[ count ].Content = []byte( partContent[ id3pos: ] )
			r.parts[ count ].MultpartType = PARTTYPE_APPLICATION_OCTET_STREAM
			r.octetStreamParts += 1


		} else {
			ctypepos := bytes.Index(partContent, []byte("Content-Type:"))
			cjsonpos := bytes.Index(partContent, []byte("application/json; charset=UTF-8")) + len(string("application/json; charset=UTF-8"))
			cjson := bytes.Index(partContent, []byte("{"))

			r.parts[ count ].ContentType = string( partContent[ ctypepos: ctypepos + cjsonpos ] )
			r.parts[ count ].Content = []byte( partContent[ cjson: ] )
			r.parts[ count ].MultpartType = PARTTYPE_APPLICATION_JSON
			r.jsonParts += 1
		}
	}
}

func ( r * DirectiveReader ) ReadAll() []ResponseDirective {

	r.analyseContent()

	directives := make( []ResponseDirective , r.jsonParts )
	directiveCount := 0

	for _ , jsonPart := range r.parts {
		if jsonPart.MultpartType == PARTTYPE_APPLICATION_JSON {
			json.Unmarshal( jsonPart.Content , &directives[ directiveCount ] )

			if r.octetStreamParts == 0 {
				directiveCount++
				continue
			}

			for _ , octetPart := range r.parts {
				if octetPart.MultpartType == PARTTYPE_APPLICATION_OCTET_STREAM {
					cidBeginPos := strings.Index( octetPart.ContentID , "<") +1
					cidEndPos := strings.Index( octetPart.ContentID , ">")

					if cid := "cid:" + octetPart.ContentID[ cidBeginPos : cidEndPos ]; cid == (&directives[ directiveCount ]).Directive.Payload.Url {
						(&directives[ directiveCount ]).Directive.mp3data = octetPart.Content

					}
				}
			}

			directiveCount++
		}
	}

	return directives
}
