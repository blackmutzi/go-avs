package main

import (
	"auth"
	"fmt"
	"http2"
	"event"
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

	go client.CreateDownchannel()

	client.Do( event.NewSystemRequest() )
	time.Sleep( 2000 * time.Millisecond )

	client.Do( event.NewSettingsRequest("de-DE" ) )
	time.Sleep( 3000 * time.Millisecond )

	client.Do( event.NewSpeechRecognizeWakeWordRequest( []int16{} ))
	time.Sleep( 3000 * time.Millisecond )

	for {
		time.Sleep( 1000 * time.Millisecond )
	}

}
