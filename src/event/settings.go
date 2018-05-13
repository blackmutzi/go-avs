package event

import (
	"io/ioutil"
	"strings"
)

type Settings struct {
	MessageID string
}

/*
	 Accepted Keys: locale
	 Accepted Values: en-AU, en-CA, en-GB, en-IN, en-US, de-DE, ja-JP
	 for more info: https://developer.amazon.com/de/docs/alexa-voice-service/settings.html
 */
func ( s * Settings ) CreateSettingsUpdateEvent( key string , value string ) string {
	var content string

	bytes , _ := ioutil.ReadFile( ASSET_PATH + "SettingsUpdateEvent.json" )
	content = string( bytes )

	content = strings.Replace( content , "{{MESSAGE_ID_STRING}}" , s.MessageID , -1 )
	content = strings.Replace( content , "{{KEY_STRING}}" , key , -1 )
	content = strings.Replace( content , "{{VALUE_STRING}}" , value  , -1 )

	return content
}

