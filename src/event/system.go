package event

import (
	"io/ioutil"
	"strings"
	"http2"
)

type System struct {
	Event * SynchronizeStateEvent
	MessageID string
}

func ( s * System ) CreateSynchronizeStateEvent() string {
	var content string

	bytes, _ := ioutil.ReadFile( ASSET_PATH + "SystemSynchronize.json")
	content = string( bytes )

	content = strings.Replace( content , "{{Alerts.AlertsState}}" , s.Event.GetAlertContext() , -1 )
	content = strings.Replace( content , "{{AudioPlayer.PlaybackState}}" , s.Event.GetPlaybackContext() , -1 )
	content = strings.Replace( content , "{{Speaker.VolumeState}}" , s.Event.GetSpeakerContext() , -1 )
	content = strings.Replace( content , "{{SpeechSynthesizer.SpeechState}}" , s.Event.GetSpeechSynthesizerContext() , -1 )
	content = strings.Replace( content , "{{SpeechRecognizer.RecognizerState}}" , s.Event.GetSpeechRecognizerContext() , -1 )
	content = strings.Replace( content , "{{MESSAGE_ID_STRING}}" , s.MessageID , -1 )

	return content
}

/*
	Create a new system request

	Reference:
 	https://developer.amazon.com/de/docs/alexa-voice-service/system.html
 */
func NewSystemRequest() * http2.Request {
	system := &System{}
	system.Event = NewSyncStateEvent()
	system.MessageID = NewMessageId()

	syncInfo := NewTransportInfo("1390402302040")
	req := &http2.Request{}
	req .TransportInfo = syncInfo.CreateMessage( system.CreateSynchronizeStateEvent() )
	return req
}

func UpdateSystemRequest( sync * SynchronizeStateEvent ) * http2.Request {
	system := &System{}
	system.Event = sync
	system.MessageID = NewMessageId()

	syncInfo := NewTransportInfo("1390402302040")
	req := &http2.Request{}
	req .TransportInfo = syncInfo.CreateMessage( system.CreateSynchronizeStateEvent() )
	return req
}





