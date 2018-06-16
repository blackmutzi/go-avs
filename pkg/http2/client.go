package http2

import (
	"net/http"
	"io/ioutil"
	"github.com/blackmutzi/go-avs/pkg/event"
	"fmt"
	"encoding/json"
	"bytes"
)

type Request struct {
	TransportInfo *event.TransportInfo
}

type ExceptionHeader struct {
	Namespace string `json:"namespace"`
	Type string `json:"name"`
	MessageID string `json:"messageId"`
}

type ExceptionPayload struct {
	Code string `json:"code"`
	Description string `json:"description"`
}

type Exception struct {
	Header ExceptionHeader `json:"header"`
	Payload ExceptionPayload `json:"payload"`
}

type Client struct {
	AccessToken string
	Version string
	DirectivesPath string
	EventsPath string
	EndPointURL string
}

/*
	 create a new http2 client for Alexa Voice Service
	 url: ASIA_ENDPOINT_URL , EU_ENDPOINT_URL , NA_ENDPOINT_URL
	 token: access token
     version: VERSION - currently ( v20160207 )
 */
func NewClient( url string , token string , version string ) *Client {
	client := &Client{}
	client.AccessToken = token
	client.Version = version
	client.DirectivesPath = "/" + version + "/directives"
	client.EventsPath = "/" + version + "/events"
	client.EndPointURL = url
	return client
}

func ( c * Client ) checkStatusCode( resp *http.Response ) (more bool, err error) {
	switch resp.StatusCode {
	case 200:
		// Keep going.
		fmt.Println("Status: 200")
		return true, nil
	case 204:
		// No content.
		fmt.Println("Status: 204")
		return false, nil
	default:
		// Attempt to parse the response as a System.Exception message.
		var exception Exception
		data, _ := ioutil.ReadAll(resp.Body)
		json.Unmarshal(data, &exception)
		if exception.Payload.Code != "" {

			fmt.Printf("Exception by: %s (%s) \n" , exception.Header.Namespace , exception.Header.MessageID )
			fmt.Printf("Code: %s \n" , exception.Payload.Code )
			fmt.Printf("Description: %s \n" , exception.Payload.Description )

			return false, fmt.Errorf("request failed with %s", resp.Status)
		}
		// Fallback error.
		return false, fmt.Errorf("request failed with %s", resp.Status)
	}
}

func ( c * Client ) CreateDownchannel() ( err error ) {
	req , err := http.NewRequest("GET", c.EndPointURL + c.DirectivesPath  , nil )
	req.Header.Add("authorization", "Bearer " + c.AccessToken )
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}

	c.checkStatusCode( resp )

	defer resp.Body.Close()
	bytes , _ := ioutil.ReadAll( resp.Body )

	fmt.Println("Downchannel Response:")
	fmt.Println( string( bytes ) )
	fmt.Println("Downchannel Bye, Bye!")
	return err
}

func ( c * Client ) Do( request *Request )( response []byte , err error ){
	req , err := http.NewRequest("POST", c.EndPointURL + c.EventsPath , bytes.NewReader( request.TransportInfo.Message.Bytes() ) )
	req.Header.Add("authorization", "Bearer " + c.AccessToken )
	req.Header.Add("content-type", "multipart/form-data; boundary=" + request.TransportInfo.Boundary )

	fmt.Println( string( request.TransportInfo.Message.Bytes() ) )

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil , err
	}
	if status , _ := c.checkStatusCode( resp ); status {
		defer resp.Body.Close()
		response , err = ioutil.ReadAll( resp.Body )
	} else {
		// no content
		fmt.Println("Request: no content ")
		response = nil
	}

	return response , err
}

/*
	Create a new Settings Request
	Accepted Values: en-AU, en-CA, en-GB, en-IN, en-US, de-DE, ja-JP
 */
func NewSettingsRequest( acceptedValue string )( *Request ) {
	settings := &event.Settings{}
	settings.MessageID = event.NewMessageId()

	settingsInfo := event.NewTransportInfo("1390402302040")
	req := &Request{}
	req.TransportInfo = settingsInfo.CreateMessage( settings.CreateSettingsUpdateEvent("locale", acceptedValue ) )
	return req
}

/*
	create a new SpeechRecognize Request with the WakeWord Profil
 */
func NewSpeechRecognizeWakeWordRequest( pcmBytes []int16 ) ( *Request ) {
	recognize := event.NewSpeechRecognizeWakeWordProfil( event.NewSyncStateEvent() )
	recogInfo := event.NewTransportInfo("1390402302040")
	req := &Request{}
	req.TransportInfo = recogInfo.CreateMessageWithAudioContent( recognize.CreateSpeechRecognizeEvent() , recogInfo.CreateAudio( pcmBytes ) )
	return req
}

/*
	create a new SpeechRecognize Request with the WakeWord Profil Default
 */
func NewSpeechRecognizeWakeWordRequestDefault( pcmLittleEndianBytes []byte ) ( *Request ) {
	recognize := event.NewSpeechRecognizeWakeWordProfil( event.NewSyncStateEvent() )
	recogInfo := event.NewTransportInfo("1390402302040")
	req := &Request{}
	req.TransportInfo = recogInfo.CreateMessageWithAudioContent( recognize.CreateSpeechRecognizeEvent() , pcmLittleEndianBytes )
	return req
}

/*
	create a new SpeechRecognize Request with the TAP Profil
 */
func NewSpeechRecognizeTAPRequest( pcmBytes []int16 ) ( *Request ) {
	recognize := event.NewSpeechRecognizeTAPProfil( event.NewSyncStateEvent() )
	recogInfo := event.NewTransportInfo("1390402302040")
	req := &Request{}
	req.TransportInfo = recogInfo.CreateMessageWithAudioContent( recognize.CreateSpeechRecognizeEvent() , recogInfo.CreateAudio( pcmBytes ) )
	return req
}

/*
	Create a new system request

	Reference:
 	https://developer.amazon.com/de/docs/alexa-voice-service/system.html
 */
func NewSystemRequest() * Request {
	system := &event.System{}
	system.Event = event.NewSyncStateEvent()
	system.MessageID = event.NewMessageId()

	syncInfo := event.NewTransportInfo("1390402302040")
	req := &Request{}
	req .TransportInfo = syncInfo.CreateMessage( system.CreateSynchronizeStateEvent() )
	return req
}

func UpdateSystemRequest( sync * event.SynchronizeStateEvent ) * Request {
	system := &event.System{}
	system.Event = sync
	system.MessageID = event.NewMessageId()

	syncInfo := event.NewTransportInfo("1390402302040")
	req := &Request{}
	req .TransportInfo = syncInfo.CreateMessage( system.CreateSynchronizeStateEvent() )
	return req
}




