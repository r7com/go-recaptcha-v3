go-recaptcha
============

About
-----

This package handles [reCaptcha](https://www.google.com/recaptcha) ([API version 3.0](https://developers.google.com/recaptcha/intro)) form submissions in [Go](http://golang.org/).

Usage
-----

Install the package in your environment:

```
go get github.com/r7com/go-recaptcha
```

To use it within your own code, import <tt>github.com/r7com/go-recaptcha</tt> and call:

```
recaptcha.Init(recaptchaPrivateKey, recaptchaScore)
```

once, to set the reCaptcha private key for your domain, then:

```
recaptcha.Confirm(recaptchaResponse)
```

for each reCaptcha form input you need to check, using the values obtained by reading the form's POST parameters (the <tt>recaptchaResponse</tt> in the above corresponds to the value of <tt>g-recaptcha-response</tt> sent by the reCaptcha server.)

The recaptcha.Confirm() function returns **true** if the captcha was completed correctly and the score was equal or above the value passed in reCaptcha.Init() or **false** if the captcha had an invalid token or the score failed, along with any errors (from the HTTP io read or the attempt to unmarshal the JSON reply).

Usage Example
-------------

Included with this repo is [example.go](example/example.go), a simple HTTP server which creates the reCaptcha form and tests the input.

See the [instructions](example/README.md) for running the example for more details.

Disclaimer
-------------
This project was forked from [dpapathanasiou](github.com/dpapathanasiou/go-recaptcha) github repository and modified to reCaptcha V3 requests.
