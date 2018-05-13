package event

import (
	"io/ioutil"
	"strings"
	"strconv"
)

const (
	ASSET_PATH = "/asset/"
)

type AlertState struct {
	Token string
	Type string
	ScheduledTime string
}

type VolumeState struct {
	Volume int
	Muted string
}

type SpeechState struct {
	Token string
	OffsetInMillisesonds int
	PlayerActivity string
}

type PlaybackState struct {
	Token string
	OffsetInMilliseconds int
	PlayerActivity string
}

type RecognizeState struct {
	Wakeword string
}

type SynchronizeStateEvent struct {
	AlertsState    *AlertState
	PlaybackState  *PlaybackState
	VolumeState    *VolumeState
	SpeechState    *SpeechState
	RecognizeState *RecognizeState
}

func NewSyncStateEvent() *SynchronizeStateEvent {
	sync := &SynchronizeStateEvent{}

	playback := &PlaybackState{}
	playback.Token = ""
	playback.OffsetInMilliseconds = 0
	playback.PlayerActivity = "IDLE"

	speechState := &SpeechState{}
	speechState.Token = ""
	speechState.OffsetInMillisesonds = 0
	speechState.PlayerActivity = "FINISHED"

	recognizeState := &RecognizeState{}
	recognizeState.Wakeword = "ALEXA"

	volumeState := &VolumeState{}
	volumeState.Muted = "false"
	volumeState.Volume = 100

	sync.PlaybackState = playback
	sync.AlertsState = &AlertState{}
	sync.SpeechState = speechState
	sync.RecognizeState = recognizeState
	sync.VolumeState = volumeState

	return sync
}

func ( s * SynchronizeStateEvent ) GetAlertContext() string {
	bytes , _ := ioutil.ReadFile( ASSET_PATH + "AlertContext.json")
	return string( bytes )
}

func (s * SynchronizeStateEvent ) GetPlaybackContext() string {
	var content string

	bytes , _ := ioutil.ReadFile( ASSET_PATH + "PlaybackContext.json" )
	content = string( bytes )

	content = strings.Replace( content , "{{TOKEN_STRING}}" , s.PlaybackState.Token , -1 )
	content = strings.Replace( content , "{{MS}}", strconv.Itoa( s.PlaybackState.OffsetInMilliseconds ) , -1)
	content = strings.Replace( content , "{{ACTIVITY_STRING}}", s.PlaybackState.PlayerActivity , -1 )

	return content
}

func ( s * SynchronizeStateEvent ) GetSpeakerContext() string {
	var content string

	bytes , _ := ioutil.ReadFile( ASSET_PATH + "SpeakerContext.json" )
	content = string( bytes )

	content = strings.Replace( content , "{{{VOLUME}}" , strconv.Itoa( s.VolumeState.Volume ) , -1 )
	content = strings.Replace( content , "{{MUTED}}" , s.VolumeState.Muted , -1 )

	return content
}

func ( s * SynchronizeStateEvent ) GetSpeechRecognizerContext() string {
	var content string

	bytes , _ := ioutil.ReadFile( ASSET_PATH + "SpeechRecognizeContext.json" )
	content = string( bytes )

	content = strings.Replace( content , "{{WAKEWORD_STRING}}" , s.RecognizeState.Wakeword , -1 )

	return content
}

func ( s * SynchronizeStateEvent ) GetSpeechSynthesizerContext() string {
	var content string

	bytes , _ := ioutil.ReadFile( ASSET_PATH + "SpeechSynthesizerContext.json" )
	content = string( bytes )

	content = strings.Replace( content , "{{TOKEN_STRING}}" , s.SpeechState.Token , -1 )
	content = strings.Replace( content , "{{MS}}" , strconv.Itoa( s.SpeechState.OffsetInMillisesonds ) , -1 )
	content = strings.Replace( content , "{{ACTIVITY_STRING}}" , s.SpeechState.PlayerActivity , -1 )

	return content
}

