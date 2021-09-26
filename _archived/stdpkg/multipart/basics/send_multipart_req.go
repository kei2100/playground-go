package main

import (
	"bytes"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"net/http/httputil"
)

func simple() {
	var body bytes.Buffer

	mw := multipart.NewWriter(&body)
	fw, err := mw.CreateFormField("medium[type]")
	if err != nil {
		panic(err)
	}
	if _, err := io.WriteString(fw, "text/csv"); err != nil {
		panic(err)
	}

	fw, err = mw.CreateFormFile("medium[file]", "test.csv")
	if err != nil {
		panic(err)
	}
	if _, err := io.WriteString(fw, "test,test"); err != nil {
		panic(err)
	}

	if err := mw.Close(); err != nil {
		panic(err)
	}
	res, err := http.Post("http://localhost:5000/media", mw.FormDataContentType(), &body)
	if err != nil {
		panic(err)
	}
	if err != nil {
		panic(err)
	}
	defer res.Body.Close()
	dump, e := httputil.DumpResponse(res, true)
	log.Println(string(dump), e)
}

func pipe() {
	pr, pw := io.Pipe()
	mw := multipart.NewWriter(pw)

	go func() {
		defer pw.Close()
		defer mw.Close()

		fw, err := mw.CreateFormField("medium[type]")
		if err != nil {
			log.Println(err)
			return
		}
		if _, err := io.WriteString(fw, "text/csv"); err != nil {
			log.Println(err)
			return
		}

		fw, err = mw.CreateFormFile("medium[file]", "test.csv")
		if err != nil {
			log.Println(err)
			return
		}
		if _, err := io.WriteString(fw, "test,test"); err != nil {
			log.Println(err)
			return
		}
	}()

	res, err := http.Post("http://localhost:5000/media", mw.FormDataContentType(), pr)
	if err != nil {
		panic(err)
	}
	if err != nil {
		panic(err)
	}
	defer res.Body.Close()
	dump, e := httputil.DumpResponse(res, true)
	log.Println(string(dump), e)
}
