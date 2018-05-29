package main

import (
	"auth"
	"fmt"
	"http2"
	"time"
	"event"
	"directive"
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

	client.Do( http2.NewSystemRequest() )
	time.Sleep( 2000 * time.Millisecond )

	client.Do( http2.NewSettingsRequest("de-DE" ) )
	time.Sleep( 3000 * time.Millisecond )

	response , err := client.Do( http2.NewSpeechRecognizeWakeWordRequest( event.ReadPCMFile("alexa_guten_morgen.wav") ) )
	fmt.Println( string( response ) )

	for _ , directive := range directive.NewDirectiveReader( response , "--------abcde123") {
		if directive.Header.Namespace == "SpeechSynthesizer" && directive.Header.Name == "Speak" {
			if directive.HasMP3Data() {

				// Play Sound
				// go playMP3Sound( directive.GetMP3Data() )
				fmt.Println("Play MP3 Sound now ...")
			}
		}
	}


	for {
		time.Sleep( 1000 * time.Millisecond )
	}

}
