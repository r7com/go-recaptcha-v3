// example.go
//
// A simple HTTP server which presents a reCaptcha input form and evaulates the result,
// using the github.com/dpapathanasiou/go-recaptcha package.
//
// See the main() function for usage.
package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"

	recaptcha "gitlab.ir7.com.br/r7/go-recaptcha"
)

var recaptchaPublicKey string

const (
	pageTop = `<!DOCTYPE HTML><html><head>
<style>.error{color:#ff0000;} .ack{color:#0000ff;}</style><title>Recaptcha Test</title></head>
<body><div style="width:100%"><div style="width: 50%;margin: 0 auto;">
<h3>Recaptcha Test</h3>
<p>This is a token generator</p>`
	form = `
    <!-- %s -->
    <script src='https://www.google.com/recaptcha/api.js?render=6LdomoYUAAAAALNFVBV_Etyu_v72PtvegyLTdm03'></script>

    <script>
      grecaptcha.ready(function() {
        grecaptcha.execute('6LdomoYUAAAAALNFVBV_Etyu_v72PtvegyLTdm03', {action: 'action_name'})
        .then(function(token) {
          document.getElementById("g-recaptcha-response-holder").innerHTML= token
        });
      });

      function generateToken() {
        grecaptcha.execute('6LdomoYUAAAAALNFVBV_Etyu_v72PtvegyLTdm03', {action: 'action_name'})
          .then(function(token) {
            document.getElementById("g-recaptcha-response-holder").innerHTML= token
        });
      }
    </script>
    <textarea id="g-recaptcha-response-holder" style="margin: 0px; width: 600px; height: 100px;"></textarea>
    </br>
    <button id="g-recaptcha-generate" onclick="generateToken()">Generate new token</button>`
	// form = `
	//   <script src="https://sc.r7.com/r7/captcha/v3/r7-recaptcha-bundle.js"></script>
	//   <script>
	//     // %s
	//     window.onload = function() {
	//       R7Recaptcha.injectRecaptcha('6LdomoYUAAAAALNFVBV_Etyu_v72PtvegyLTdm03')
	//     }
	//
	//     function generateToken() {
	//       R7Recaptcha.newToken((token)=>{ document.getElementById("g-recaptcha-response-holder").innerHTML= token }, "homepage")
	//     }
	//   </script>
	//   <textarea id="g-recaptcha-response-holder" style="margin: 0px; width: 600px; height: 100px;"></textarea>
	//   </br>
	//   <button id="g-recaptcha-generate" onclick="generateToken()">Generate new token</button>`
	pageBottom = `</div></div></body></html>`
	anError    = `<p class="error">%s</p>`
	anAck      = `<p class="ack">%s</p>`
)

// processRequest accepts the http.Request object, finds the reCaptcha form variables which
// were input and sent by HTTP POST to the server, then calls the recaptcha package's Confirm()
// method, which returns a boolean indicating whether or not the client answered the form correctly.
func processRequest(request *http.Request) (result bool) {
	recaptchaResponse, responseFound := request.Form["g-recaptcha-response"]
	if responseFound {
		result, err := recaptcha.Confirm("127.0.0.1", recaptchaResponse[0])
		if err != nil {
			log.Println("recaptcha server error", err)
		}
		return result
	}
	return false
}

// homePage is a simple HTTP handler which produces a basic HTML page
// (as defined by the pageTop and pageBottom constants), including
// an input form with a reCaptcha challenge.
// If the http.Request object indicates the form input has been posted,
// it calls processRequest() and displays a message indicating whether or not
// the reCaptcha form was input correctly.
// Either way, it writes HTML output through the http.ResponseWriter.
func homePage(writer http.ResponseWriter, request *http.Request) {
	err := request.ParseForm() // Must be called before writing response
	fmt.Fprint(writer, pageTop)
	if err != nil {
		fmt.Fprintf(writer, fmt.Sprintf(anError, err))
	} else {
		_, buttonClicked := request.Form["button"]
		if buttonClicked {
			if processRequest(request) {
				// fmt.Fprint(writer, fmt.Sprintf(anAck, "Recaptcha was correct!"))
			} else {
				// fmt.Fprintf(writer, fmt.Sprintf(anError, "Recaptcha was incorrect; try again."))
			}
		}
	}
	fmt.Fprint(writer, fmt.Sprintf(form, recaptchaPublicKey))
	fmt.Fprint(writer, pageBottom)
}

// main expects two command-line arguments: the reCaptcha public key for producing the HTML form,
// and the reCaptcha private key, to pass to recaptcha.Init() so the recaptcha package can check the input.
// It launches a simple web server on port 9001 which produces the reCaptcha input form and checks the client
// input if the form is posted.
func main() {
	if len(os.Args) != 3 {
		fmt.Printf("usage: %s <reCaptcha public key> <reCaptcha private key>\n", filepath.Base(os.Args[0]))
		os.Exit(1)
	} else {
		recaptchaPublicKey = os.Args[1]
		recaptcha.Init(os.Args[2])

		http.HandleFunc("/", homePage)
		if err := http.ListenAndServe(":9001", nil); err != nil {
			log.Fatal("failed to start server", err)
		}
	}
}
