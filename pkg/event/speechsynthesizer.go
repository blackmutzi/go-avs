package event

import (
	"io/ioutil"
	"strings"
)

type SpeechSynthesizer struct {
	SpeechState *SpeechState
	MessageID string
}

func ( s * SpeechSynthesizer ) createSpeechStartedEvent() string {
	var content string

	bytes , _ := ioutil.ReadFile( ASSET_PATH + "SpeechStartedEvent.json" )
	content = string( bytes )

	content = strings.Replace( content , "{{MESSAGE_ID_STRING}}" , s.MessageID , -1 )
	content = strings.Replace( content , "{{TOKEN_STRING}}" , s.SpeechState.Token , -1 )

	return content
}

func ( s * SpeechSynthesizer ) createSpeechFinishedEvent() string {
	var content string

	bytes , _ := ioutil.ReadFile( ASSET_PATH + "SpeechFinishedEvent.json" )
	content = string( bytes )

	content = strings.Replace( content , "{{MESSAGE_ID_STRING}}" , s.MessageID , -1 )
	content = strings.Replace( content , "{{TOKEN_STRING}}" , s.SpeechState.Token , -1 )

	return content
}


