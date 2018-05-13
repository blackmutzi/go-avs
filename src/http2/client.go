package http2

import (
	"net/http"
	"io/ioutil"
	//"encoding/json"
	"event"
	"fmt"
	"strings"
)

type Request struct {
	TransportInfo *event.TransportInfo
}

type Response struct {

}

type Exception struct {

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
		return true, nil
	case 204:
		// No content.
		return false, nil
	default:
		// Attempt to parse the response as a System.Exception message.
		data, _ := ioutil.ReadAll(resp.Body)
		fmt.Println("Exception:")
		fmt.Println( data )

		//var exception Exception
		//json.Unmarshal(data, &exception)
		//if exception.Payload.Code != "" {
		//	return false, &exception
		//}
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
	if status , err := c.checkStatusCode( resp ); !status {
		return err
	}
	defer resp.Body.Close()
	bytes , _ := ioutil.ReadAll( resp.Body )
	fmt.Println( bytes )

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
	if status , err := c.checkStatusCode( resp ); !status {
		return err
	}
	defer resp.Body.Close()
	bytes , _ := ioutil.ReadAll( resp.Body )
	fmt.Println( bytes )

	return err
}



