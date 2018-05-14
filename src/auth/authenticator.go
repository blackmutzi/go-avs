package auth

import (
	"encoding/json"
	"io/ioutil"
	"strings"
	"net/http"
	"errors"
	"net/url"
)

type AuthInfo struct {
	ClientID string `json:"client_id"`
	ClientSecret string `json:"client_secret"`
	ProductName string `json:"product_id"`
	RedirectURI string `json:"redirect_uri"`
	GrantCode string `json:"code_grant"`
	RefreshToken string `json:"refresh_token"`
	AccessToken string `json:"access_token"`
}

type Authenticator struct {
	OAuth2URL string
	AuthInfo *AuthInfo
}

func NewAuth( jsonConfigFile string ) *Authenticator {
	auth := &Authenticator{}
	auth.OAuth2URL = "https://api.amazon.com/auth/o2/token"
	auth.AuthInfo = &AuthInfo{}
	bytes , _ := ioutil.ReadFile( jsonConfigFile )
	json.Unmarshal( bytes , auth.AuthInfo )
	return auth
}

func ( auth * Authenticator) WriteFile( fileName string ) ( err error ){
	data , _ := json.Marshal( auth.AuthInfo )
	err = ioutil.WriteFile( fileName , data , 0644)
	return err
}

func ( auth * Authenticator ) CreateAmazonLoginLink() string {
	var amazonLink string
	scope := "alexa:all"
	scopeData := "{\"${SCOPE}\": {\"productID\": \"${PRODUCT_ID}\",\"productInstanceAttributes\": {\"deviceSerialNumber\": \"${DEVICE_SERIAL_NUMBER}\"}}}"

	scopeData = strings.Replace( scopeData , "${SCOPE}" , scope , 1 )
	scopeData = strings.Replace( scopeData , "${PRODUCT_ID}" , auth.AuthInfo.ProductName , 1 )
	scopeData = strings.Replace( scopeData , "${DEVICE_SERIAL_NUMBER}" , "12345" , 1 )

	amazonLink += "client_id=" + auth.AuthInfo.ClientID
	amazonLink += "&scope=" + scope
	amazonLink += "&scope_data=" + scopeData
	amazonLink += "&response_type=code"
	amazonLink += "&redirect_uri=" + auth.AuthInfo.RedirectURI

	return "https://www.amazon.com/ap/oa?" + amazonLink
}

func ( auth * Authenticator ) getRefreshToken() ( err error ){
	var authTemp AuthInfo

	formData := url.Values{}
	formData.Set("grant_type", "authorization_code" )
	formData.Set("code" , auth.AuthInfo.GrantCode )
	formData.Set("client_id", auth.AuthInfo.ClientID )
	formData.Set("client_secret" , auth.AuthInfo.ClientSecret )
	formData.Set("redirect_uri", auth.AuthInfo.RedirectURI )

	req, err := http.NewRequest( "POST" , auth.OAuth2URL , strings.NewReader( formData.Encode()) )
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	resp, err := http.DefaultClient.Do( req )
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	bytes , err := ioutil.ReadAll( resp.Body )
	if err != nil {
		return err
	}

	json.Unmarshal( bytes  , &authTemp )

	if len( authTemp.RefreshToken ) != 0 {
		auth.AuthInfo.RefreshToken = authTemp.RefreshToken
	} else {
		err = errors.New("invalid grant code or Bad Request")
	}
	return err
}

func ( auth * Authenticator ) GetAccessToken() ( err error ){
	var authTemp AuthInfo

	if len( auth.AuthInfo.RefreshToken ) == 0 {
		err = auth.getRefreshToken()
		if err != nil {
			return err
		}
	}

	formData := url.Values{}
	formData.Set("grant_type", "refresh_token" )
	formData.Set("refresh_token" , auth.AuthInfo.RefreshToken  )
	formData.Set("client_id", auth.AuthInfo.ClientID )
	formData.Set("client_secret" , auth.AuthInfo.ClientSecret )

	req, err := http.NewRequest( "POST" , auth.OAuth2URL , strings.NewReader( formData.Encode() ))
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	resp, err := http.DefaultClient.Do( req )
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	bytes , err := ioutil.ReadAll( resp.Body )
	if err != nil {
		return err
	}
	json.Unmarshal( bytes  , &authTemp )

	if len( authTemp.RefreshToken ) != 0 {
		auth.AuthInfo.RefreshToken = authTemp.RefreshToken
	}

	if len( authTemp.AccessToken ) != 0 {
		auth.AuthInfo.AccessToken = authTemp.AccessToken
	} else {
		err = errors.New("bad request")
	}

	return err
}

