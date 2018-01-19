package main

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"time"

	fdk "github.com/fnproject/fdk-go"
)

type AppRequest struct {
	// if specified we 'sleep' the specified msecs
	SleepTime int `json:"sleepTime,omitempty"`
	// if specified, this is our response http status code
	ResponseCode int `json:"responseCode,omitempty"`
	// if specified, this is our response content-type
	ResponseContentType string `json:"responseContentType,omitempty"`
	// if specified, this is echoed back to client
	EchoContent string `json:"echoContent,omitempty"`
	// verbose mode
	IsDebug bool `json:"isDebug,omitempty"`
	// simulate crash
	IsCrash bool `json:"isCrash,omitempty"`
	// TODO: simulate slow read/slow write
	// TODO: simulate partial write/read
	// TODO: simulate mem leak
	// TODO: simulate high cpu usage
	// TODO: simulate high mem usage
	// TODO: simulate large body upload/download
}

type AppResponse struct {
	Request AppRequest        `json:"request"`
	Headers http.Header       `json:"header"`
	Config  map[string]string `json:"config"`
}

func AppHandler(ctx context.Context, in io.Reader, out io.Writer) {
	fnctx := fdk.Context(ctx)

	var b bytes.Buffer
	io.Copy(&b, in)
	in = bytes.NewReader(b.Bytes())

	var request AppRequest
	json.NewDecoder(in).Decode(&request)

	if request.IsDebug {
		log.Printf("Received request %v", request)
		log.Printf("Received headers %v", fnctx.Header)
		log.Printf("Received config %v", fnctx.Config)
	}

	// simulate load if requested
	if request.SleepTime > 0 {
		if request.IsDebug {
			log.Printf("Sleeping %d", request.SleepTime)
		}
		time.Sleep(time.Duration(request.SleepTime) * time.Millisecond)
	}

	// simulate crash
	if request.IsCrash {
		panic("Crash requested")
	}

	// custom content type
	if request.ResponseContentType != "" {
		fdk.SetHeader(out, "Content-Type", request.ResponseContentType)
	} else {
		fdk.SetHeader(out, "Content-Type", "application/json")
	}

	// custom response code
	if request.ResponseCode != 0 {
		fdk.WriteStatus(out, request.ResponseCode)
	} else {
		fdk.WriteStatus(out, 200)
	}

	resp := AppResponse{
		Request: request,
		Headers: fnctx.Header,
		Config:  fnctx.Config,
	}

	io.Copy(out, &b)

	// json.NewEncoder(out).Encode(&resp)
	_ = resp
}

func main() {
	fdk.Handle(fdk.HandlerFunc(AppHandler))
}
