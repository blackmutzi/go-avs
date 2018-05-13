package event

import (
	"io/ioutil"
	"strings"
)

type SpeechRecognize struct {
	Event *SynchronizeStateEvent
	MessageID string
	DialogRequestID string
}

func ( s * SpeechRecognize ) createSpeechRecognizeEvent() string {
	var content string

	bytes , _ := ioutil.ReadFile( ASSET_PATH + "SpeechRecognizeEvent.json" )
	content = string( bytes )

	content = strings.Replace( content , "{{Alerts.AlertsState}}" , s.Event.GetAlertContext() , -1 )
	content = strings.Replace( content , "{{AudioPlayer.PlaybackState}}" , s.Event.GetPlaybackContext() , -1 )
	content = strings.Replace( content , "{{Speaker.VolumeState}}" , s.Event.GetSpeakerContext() , -1 )
	content = strings.Replace( content , "{{SpeechSynthesizer.SpeechState}}" , s.Event.GetSpeechSynthesizerContext() , -1 )
	content = strings.Replace( content , "{{SpeechRecognizer.RecognizerState}}" , s.Event.GetSpeechRecognizerContext() , -1 )
	content = strings.Replace( content , "{{MESSAGE_ID_STRING}}" , s.MessageID , -1 )

	content = strings.Replace( content , "{{DIALOG_STRING}}" , s.DialogRequestID , -1 )
	content = strings.Replace( content , "{{PROFILE_STRING}}" , "NEAR_FIELD" , -1 )
	content = strings.Replace( content , "{{FORMAT_STRING}}" , "AUDIO_L16_RATE_16000_CHANNELS_1" , -1 )

	content = strings.Replace( content , "{{TYPE_STRING}}" , "TAP" , -1 )
	content = strings.Replace( content , "{{LONG}}" , "0" , -1 )

	return content
}
