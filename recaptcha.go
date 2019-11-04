// Package recaptcha handles reCaptcha (http://www.google.com/recaptcha) form submissions
//
// This package is designed to be called from within an HTTP server or web framework
// which offers reCaptcha form inputs and requires them to be evaluated for correctness
//
// Edit the recaptchaPrivateKey constant before building and using
package recaptcha

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"time"
)

type RecaptchaResponse struct {
	Success     bool      `json:"success"`
	Score       float32   `json:"score"`
	ChallengeTS time.Time `json:"challenge_ts"`
	Hostname    string    `json:"hostname"`
	ErrorCodes  []string  `json:"error-codes"`
}

const recaptchaServerName = "https://www.google.com/recaptcha/api/siteverify"

var recaptchaPrivateKey string
var recaptchaScore float32
var timeResponse int
var postError bool

// check uses the client ip address, the challenge code from the reCaptcha form,
// and the client's response input to that challenge to determine whether or not
// the client answered the reCaptcha input question correctly.
// It returns a boolean value indicating whether or not the client answered correctly.
func check(response string) (r RecaptchaResponse, err error) {
	postError = false

	netClient := &http.Client{
		Timeout: time.Duration(timeResponse) * time.Second,
	}

	resp, err := netClient.PostForm(recaptchaServerName,
		url.Values{"secret": {recaptchaPrivateKey}, "response": {response}})
	if err != nil {
		log.Printf("Post error: %s\n", err)
		postError = true
		return
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("Read error: could not read body: %s", err)
		return
	}

	err = json.Unmarshal(body, &r)
	if err != nil {
		fmt.Printf("Read error: got invalid JSON: %s", err)
		return
	}

	return
}

// Confirm is the public interface function.
// It calls check, which the client ip address, the challenge code from the reCaptcha form,
// and the client's response input to that challenge to determine whether or not
// the client answered the reCaptcha input question correctly.
// It returns a boolean value indicating whether or not the client answered correctly.
func Confirm(response string, ip string) (result bool, err error) {
	result = false
	resp, err := check(response)

	if resp.Success == true && resp.Score >= recaptchaScore {
		result = true
	}

	if postError == true {
		result = true
	}

	logCaptchaResult(result, resp.Score, ip)

	return
}

// Init allows the webserver or code evaluating the reCaptcha form input to set the
// reCaptcha private key (string) value, which will be different for every domain.
func Init(key string, score float32, time int) {
	recaptchaPrivateKey = key
	recaptchaScore = score
	timeResponse = time
}

func logCaptchaResult(success bool, score float32, ip string) {
	if success {
		log.Printf("[%v] Captcha: Valid token with score of %f\n", ip, score)
		return
	}

	if score > 0 {
		log.Printf("[%v] Captcha: Valid token but refused due low score(got: %f, expected: %f)", ip, score, recaptchaScore)
		return
	}

	log.Printf("[%v] Captcha: Invalid token", ip)
}
