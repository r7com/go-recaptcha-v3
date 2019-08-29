package recaptcha

import (
	"net/http"
	"strconv"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	httpmock "gopkg.in/jarcoal/httpmock.v1"
)

func TestConfirm(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	tests := []struct {
		score               string
		httpResponseStatus  int
		httpResponseMessage string
		expectedResult      bool
		errorMessage        string
	}{
		{"0.5", 200, `{"success": true, "score": 0.9}`, true, "It must be true"},
		{"0.5", 200, `{"success": false, "score": 0.2}`, false, "It must be false"},
		{"0.5", 200, `{"success": false}`, false, "It must be false"},
		{"0.5", 500, `{"success": false}`, false, "It must be false"},
	}

	for _, test := range tests {
		httpmock.RegisterResponder("POST", recaptchaServerName,
			httpmock.NewStringResponder(test.httpResponseStatus, test.httpResponseMessage))

		score, _ := strconv.ParseFloat(test.score, 32)
		Init("SOME_KEY", float32(score), 2)
		result, _ := Confirm("test")

		assert.Equal(t, test.expectedResult, result, test.errorMessage)
	}
}

func TestConfirmSlowResponse(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	httpmock.RegisterResponder("POST", recaptchaServerName,
		func(req *http.Request) (*http.Response, error) {
			time.Sleep(90 * time.Second)
			return httpmock.NewJsonResponse(200, map[string]interface{}{
				"success": true,
				"score":   0.9,
			})
		},
	)

	score, _ := strconv.ParseFloat("0.5", 32)
	Init("SOME_KEY", float32(score), 2)
	result, _ := Confirm("test")

	assert.Equal(t, true, result, "Timeout expired!")
}
