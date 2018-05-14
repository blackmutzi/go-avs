package http2

import (
	"net/http"
	"io/ioutil"
	//"encoding/json"
	"event"
	"fmt"
	"strings"
	"encoding/json"
)

type Request struct {
	TransportInfo *event.TransportInfo
}

type Response struct {

}

type ExceptionHeader struct {
	Namespace string `json:"namespace"`
	Type string `json:"name"`
	MessageID string `json:"messageId"`
}

type ExceptionPayload struct {
	Code string `json:"code"`
	Description string `json:"description"`
}

type Exception struct {
	Header ExceptionHeader `json:"header"`
	Payload ExceptionPayload `json:"payload"`
}

type Client struct {
	AccessToken string
	Version string
	DirectivesPath string
	EventsPath string
	EndPointURL string
}

/*
	 create a new http2 client for Alexa Voice Service
	 url: ASIA_ENDPOINT_URL , EU_ENDPOINT_URL , NA_ENDPOINT_URL
	 token: access token
     version: VERSION - currently ( v20160207 )
 */
func NewClient( url string , token string , version string ) *Client {
	client := &Client{}
	client.AccessToken = token
	client.Version = version
	client.DirectivesPath = "/" + version + "/directives"
	client.EventsPath = "/" + version + "/events"
	client.EndPointURL = url
	return client
}

func ( c * Client ) checkStatusCode( resp *http.Response ) (more bool, err error) {
	switch resp.StatusCode {
	case 200:
		// Keep going.
		fmt.Println("Status: 200")
		return true, nil
	case 204:
		// No content.
		fmt.Println("Status: 204")
		return false, nil
	default:
		// Attempt to parse the response as a System.Exception message.
		var exception Exception
		data, _ := ioutil.ReadAll(resp.Body)
		json.Unmarshal(data, &exception)
		if exception.Payload.Code != "" {

			fmt.Printf("Exception by: %s (%s) \n" , exception.Header.Namespace , exception.Header.MessageID )
			fmt.Printf("Code: %s \n" , exception.Payload.Code )
			fmt.Printf("Description: %s \n" , exception.Payload.Description )

			return false, fmt.Errorf("request failed with %s", resp.Status)
		}
		// Fallback error.
		return false, fmt.Errorf("request failed with %s", resp.Status)
	}
}

func ( c * Client ) CreateDownchannel() ( err error ) {
	req , err := http.NewRequest("GET", c.EndPointURL + c.DirectivesPath  , nil )
	req.Header.Add("authorization", "Bearer " + c.AccessToken )
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}

	c.checkStatusCode( resp )

	defer resp.Body.Close()
	bytes , _ := ioutil.ReadAll( resp.Body )

	fmt.Println("Downchannel Response:\n")
	fmt.Println( string( bytes ) )


	fmt.Println("Downchannel Bye, Bye!")
	return err
}

func ( c * Client ) Do( request *Request )( err error ){
	req , err := http.NewRequest("POST", c.EndPointURL + c.EventsPath , strings.NewReader(  request.TransportInfo.Message ) )
	req.Header.Add("authorization", "Bearer " + c.AccessToken )
	req.Header.Add("content-type", "multipart/form-data; boundary=" + request.TransportInfo.Boundary )

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	if status , _ := c.checkStatusCode( resp ); status {

		defer resp.Body.Close()
		bytes , _ := ioutil.ReadAll( resp.Body )

		fmt.Println("Custom Response:\n")
		fmt.Println( bytes )

	} else {
		// no content
		fmt.Println("Request: no content ")
	}


	return err
}



