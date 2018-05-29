package event

import (
	"io/ioutil"
	"strings"
	"http2"
)

type SpeechRecognize struct {
	Event *SynchronizeStateEvent
	MessageID string
	DialogRequestID string
	InitiatorType string
	StartIndexSamples string
	EndIndexSamples string
}
/*
	Speech Recognize WakeWord Profil

	Reference:
	https://developer.amazon.com/de/docs/alexa-voice-service/streaming-requirements-for-cloud-based-wake-word-verification.html
 */
func NewSpeechRecognizeWakeWordProfil( sync * SynchronizeStateEvent ) * SpeechRecognize {
	speech := &SpeechRecognize{}
	speech.Event = sync
	speech.MessageID = NewMessageId()
	speech.DialogRequestID = NewDialogReqId()
	speech.InitiatorType = "WAKEWORD"
	speech.StartIndexSamples = "8000"
	speech.EndIndexSamples = "16000"
	return speech
}

/*
	Default Speech Recognize TAP Profil
 */
func NewSpeechRecognizeTAPProfil( sync * SynchronizeStateEvent ) * SpeechRecognize {
	speech := &SpeechRecognize{}
	speech.Event = sync
	speech.MessageID = NewMessageId()
	speech.DialogRequestID = NewDialogReqId()
	speech.InitiatorType = "TAP"
	speech.StartIndexSamples = "0"
	speech.EndIndexSamples = "0"
	return speech
}

func ( s * SpeechRecognize ) CreateSpeechRecognizeEvent() string {
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

	content = strings.Replace( content , "{{TYPE_STRING}}" , s.InitiatorType , -1 )
	content = strings.Replace( content , "{{START_INDEX_SAMPLES}}" , "0" , -1 )
	content = strings.Replace( content , "{{END_INDEX_SAMPLES}}" , "0" , -1 )
	content = strings.Replace( content , "{{PAYLOAD_TOKEN}}" , "" , -1 )

	return content
}

/*
	Create a new SpeechRecognize Request with the WakeWord Profil
 */
func NewSpeechRecognizeWakeWordRequest( pcmBytes []int16 ) ( *http2.Request ) {
	recognize := NewSpeechRecognizeWakeWordProfil( NewSyncStateEvent() )
	recogInfo := NewTransportInfo("1390402302040")
	req := &http2.Request{}
	req.TransportInfo = recogInfo.CreateMessageWithAudioContent( recognize.CreateSpeechRecognizeEvent() , recogInfo.CreateAudio( pcmBytes ) )
	return req
}

/*
	Create a new SpeechRecognize Request with the WakeWord Profil
 */
func NewSpeechRecognizeTAPRequest( pcmBytes []int16 ) ( *http2.Request ) {
	recognize := NewSpeechRecognizeTAPProfil( NewSyncStateEvent() )
	recogInfo := NewTransportInfo("1390402302040")
	req := &http2.Request{}
	req.TransportInfo = recogInfo.CreateMessageWithAudioContent( recognize.CreateSpeechRecognizeEvent() , recogInfo.CreateAudio( pcmBytes ) )
	return req
}

