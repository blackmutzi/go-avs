package main

import (
	"auth"
	"fmt"
	"http2"
	"event"
	"github.com/satori/go.uuid"
	"time"
)

const (
	ASIA_ENDPOINT_URL = "https://avs-alexa-fe.amazon.com"
	EU_ENDPOINT_URL = "https://avs-alexa-eu.amazon.com"
	NA_ENDPOINT_URL = "https://avs-alexa-na.amazon.com"
	VERSION = "v20160207"
)

func main(){
	test := auth.NewAuth("auth-config.json")
	err := test.GetAccessToken()

	if err != nil {
		fmt.Println( err.Error() )
	} else {
		fmt.Println( test.AuthInfo.AccessToken  )
		fmt.Println( test.AuthInfo.RefreshToken )
		test.WriteFile("auth-config.json")
	}

	// Build Client
	client := http2.NewClient( EU_ENDPOINT_URL , test.AuthInfo.AccessToken , VERSION )
	system := event.System{}
	system.Event = event.NewSyncStateEvent()
	system.MessageID = fmt.Sprintf("%s", uuid.Must(uuid.NewV4()) )

	// make synchronize event request
	sync_info := event.NewTransportInfo("1390402302040" )
	req_sync := &http2.Request{}
	req_sync.TransportInfo = sync_info.CreateMessage( system.CreateSynchronizeStateEvent() )

	// setup settings event
	settings := event.Settings{}
	settings.MessageID = fmt.Sprintf("%s", uuid.Must(uuid.NewV4()) )

	// make settings event request
	settings_info := event.NewTransportInfo("1390402302040" )
	req_settings := &http2.Request{}
	req_settings.TransportInfo = settings_info.CreateMessage( settings.CreateSettingsUpdateEvent("locale", "de-DE") )

	recognize := event.SpeechRecognize{}
	recognize.MessageID = fmt.Sprintf("%s", uuid.Must(uuid.NewV4()) )
	recognize.Event = event.NewSyncStateEvent()
	recognize.DialogRequestID = "dialog-" + fmt.Sprintf("%s", uuid.Must(uuid.NewV4()) )

	recog_info := event.NewTransportInfo("1390402302040" )
	req := &http2.Request{}
	req.TransportInfo = recog_info.CreateMessageWithAudioContent( recognize.CreateSpeechRecognizeEvent() , []byte("audio_bytes") )

	go client.CreateDownchannel()

	fmt.Println( req_sync.TransportInfo.Message )
	client.Do( req_sync )
	time.Sleep( 2000 * time.Millisecond )

	fmt.Println( req_settings.TransportInfo.Message )
	client.Do( req_settings )
	time.Sleep( 3000 * time.Millisecond )

	fmt.Println( req.TransportInfo.Message )
	client.Do( req )

	for {
		time.Sleep( 1000 * time.Millisecond )
	}

}
