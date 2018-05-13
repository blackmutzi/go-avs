package main

import (
	"auth"
	"fmt"
	"http2"
	"event"
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
	// Build New Transport Message
	info := event.NewTransportInfo("1390402302040" )

	system := event.System{}
	system.Event = event.NewSyncStateEvent()
	system.MessageID = "UUID_HERE"
	// make synchronize event request
	req_sync := &http2.Request{}
	req_sync.TransportInfo = info.CreateMessage( system.CreateSynchronizeStateEvent() )

	// setup settings event
	settings := event.Settings{}
	settings.MessageID = "UUID_HERE"
	// make settings event request
	req_settings := &http2.Request{}
	req_settings.TransportInfo = info.CreateMessage( settings.CreateSettingsUpdateEvent("locale", "de-DE") )

	go client.CreateDownchannel()

	client.Do( req_sync )
	client.Do( req_settings )
}
