package main

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"
	"os"

	"context"

	"cloud.google.com/go/storage"
	"github.com/elgs/gostrgen"
	"github.com/julienschmidt/httprouter"
	"github.com/sirupsen/logrus"
)

type Response struct {
	Success bool
	Message string
	Name    string
}

type SimpleResponse struct {
	Success bool
	Message string
}

var peopleServed = 0
var successfulServed = 0
var unsuccessfulServed = 0
var client *storage.Client

type StatsResponse struct {
	PeopleServed             int
	SuccessfulPeopleServed   int
	UnsuccessfulPeopleServed int
}

func upload(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	logrus.Info("Handling upload request")
	if r.Method == "POST" {
		ctx := context.Background()

		r.ParseMultipartForm(10000000)
		key := r.Form.Get("key")

		if key != os.Getenv("UPLOAD_KEY") {
			json.NewEncoder(w).Encode(SimpleResponse{Success: false, Message: "Invalid upload key"})
			return
		}

		file, _, err := r.FormFile("img")

		if err != nil {
			json.NewEncoder(w).Encode(SimpleResponse{Success: false, Message: "Could not upload: " + err.Error()})
			return
		}

		defer file.Close()

		randName := randString(6)

		wc := client.Bucket(os.Getenv("BUCKET_NAME")).Object(randName + ".png").NewWriter(ctx)
		if _, err = io.Copy(wc, file); err != nil {
			json.NewEncoder(w).Encode(SimpleResponse{Success: false, Message: "Could not copy file"})
			return
		}
		if err := wc.Close(); err != nil {
			json.NewEncoder(w).Encode(SimpleResponse{Success: false, Message: "Could not close file"})
			return
		}

		json.NewEncoder(w).Encode(Response{Success: true, Message: "File uploaded", Name: randName})

	} else {
		json.NewEncoder(w).Encode(SimpleResponse{Success: false, Message: "Invalid method"})
	}

}

func serveImage(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	ctx := context.Background()

	imgName := ps.ByName("imgname")

	rc, err := client.Bucket(os.Getenv("BUCKET_NAME")).Object(imgName + ".png").NewReader(ctx)
	if err != nil {
		json.NewEncoder(w).Encode(SimpleResponse{Success: false, Message: "Unknown file"})
		return
	}
	defer rc.Close()

	if rc == nil {
		json.NewEncoder(w).Encode(SimpleResponse{Success: false, Message: "Unknown file"})
		return
	}

	data, err := ioutil.ReadAll(rc)

	if err != nil {
		json.NewEncoder(w).Encode(SimpleResponse{Success: false, Message: "Unknown read error"})
		return
	}

	logrus.Info("serving " + imgName)

	//w.Header().Set("Content-Disposition", "attachment; filename=image.png")
	// Display inline in the browser, not download
	w.Header().Set("Content-Disposition", "inline")
	w.Header().Set("Content-Type", "image/png")

	w.Write(data)
}

func main() {
	logrus.Info("starting...")

	ctx := context.Background()

	client, _ = storage.NewClient(ctx)

	r := httprouter.New()

	r.POST("/upload", upload)
	r.GET("/:imgname", serveImage)

	if err := http.ListenAndServe(":"+os.Getenv("PORT"), r); err != nil {
		logrus.Error(err.Error())
		os.Exit(1)
	}

}

func randString(n int) string {
	r, e := gostrgen.RandGen(n, gostrgen.Lower|gostrgen.Upper, "", "")
	if e != nil {
		logrus.Error("Could not generate random string")
	}
	return r
}
