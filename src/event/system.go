package event

import (
	"io/ioutil"
	"strings"
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





