package utils

import (
	"bytes"
	"crypto/tls"
	"encoding/base64"
	"encoding/json"
	"net/http"
)

const (
	url            = "https://www.smspanel.trez.ir/api/smsAPI/SendMessage"
	contentType    = "application/json"
	connection     = "keep-alive"
	userAgent      = "Mozilla/5.0 (Linux; Android 6.0; Nexus 5 Build/MRA58N) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/88.0.4324.96 Mobile Safari/537.36"
	accept         = "text/html,application/json,application/xhtml+xml,application/xml;q=0.9,image/webp,image/apng,*/*;q=0.8"
	acceptEncoding = "gzip, deflate, br"
	acceptLanguage = "en-US,en;q=0.9"
)

var (
	userPassTrez = GetEnv("TREZ_PASS", "user:password")
	phoneNumber  = GetEnv("TREZ_PHONENUMBER", "123456")
)

func createHeaderForTrez() string {
	result := base64.StdEncoding.EncodeToString([]byte(userPassTrez))
	return "Basic " + result
}

func marshallSms(obj interface{}) ([]byte, *RestErr) {
	result, err := json.Marshal(obj)
	if err != nil {
		return nil, NewInternalServerError("error while marshal a sms struct to json object")
	}
	return result, nil
}

func SendSms(message string, mobile string) *RestErr {
	//This is body of request for trez SMS service
	body := &struct {
		PhoneNumber         string   `json:"PhoneNumber"`
		Message             string   `json:"Message"`
		Mobiles             []string `json:"Mobiles"`
		UserGroupID         string   `json:"UserGroupID"`
		SendDateInTimeStamp int      `json:"SendDateInTimeStamp"`
	}{
		PhoneNumber:         phoneNumber,
		Message:             message,
		Mobiles:             []string{mobile},
		UserGroupID:         generateUniqueString(),
		SendDateInTimeStamp: 0,
	}
	jsonBodyObj, bodyErr := marshallSms(body)
	if bodyErr != nil {
		return bodyErr
	}
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonBodyObj))
	if err != nil {
		return NewInternalServerError(err.Error())
	}

	// Set header for POST request
	req.Host = "smspanel.trez.ir"
	req.Header.Set("Content-Type", contentType)
	req.Header.Set("Authorization", createHeaderForTrez())
	req.Header.Set("Connection", connection)
	req.Header.Set("User-Agent", userAgent)
	req.Header.Set("Upgrade-Insecure-Requests", "1")
	req.Header.Set("Accept", accept)
	req.Header.Set("Accept-Encoding", acceptEncoding)
	req.Header.Set("Accept-Language", acceptLanguage)

	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Transport: tr}
	resp, clientErr := client.Do(req)
	if clientErr != nil {
		return NewInternalServerError(clientErr.Error())
	}
	defer resp.Body.Close()
	return nil
}
