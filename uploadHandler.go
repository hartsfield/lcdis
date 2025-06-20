package main

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"slices"
	"strings"
	"time"
)

var stream []*item
var itemsMap map[string]*item = make(map[string]*item)

func itemView(id string) *item {
	return itemsMap[id]
}

func uploadHandler(w http.ResponseWriter, r *http.Request) {
	data := partFormData(r, w)
	stream = append([]*item{data}, stream...)

	b, err := json.Marshal(data)
	if err != nil {
		log.Println(err)
	}
	ajaxResponse(w, map[string]string{
		"success":    "true",
		"replyID":    data.ID,
		"itemString": string(b),
	})
	saveJSON()
}

func partFormData(r *http.Request, w http.ResponseWriter) *item {
	mr, err := r.MultipartReader()
	if err != nil {
		log.Println(err)
	}

	var data *item = &item{ID: genPostID(10)}

	for {
		part, err_part := mr.NextPart()
		if err_part == io.EOF {
			break
		}
		if part.FormName() == "FileElement" {
			handleFile(w, part, data)
		}
		if part.FormName() == "Title" {
			buf := new(bytes.Buffer)
			buf.ReadFrom(part)
			data.Title = buf.String()
		}
		if part.FormName() == "Article" {
			buf := new(bytes.Buffer)
			buf.ReadFrom(part)
			data.Article = buf.String()
		}
	}
	return data
}

func handleFile(w http.ResponseWriter, part *multipart.Part, data *item) {
	fileBytes, err := io.ReadAll(io.LimitReader(part, 10<<20))
	if err != nil {
		log.Println(err)
	}
	mt := http.DetectContentType(fileBytes)
	var fileExtension string
	switch mt {
	case "image/png":
		fileExtension = "png"
	case "image/jpeg":
		fileExtension = "jpg"
	case "image/gif":
		fileExtension = "gif"
	case "video/mp4":
		fileExtension = "mp4"
	case "video/webm":
		fileExtension = "webm"
	default:
		ajaxResponse(w, map[string]string{
			"success": "false",
			"replyID": "",
			"error":   "png - jpg - gif - webm - mp4 only",
		})
		return
	}
	tempFile, err := os.CreateTemp("public/temp", "u-*."+fileExtension)
	if err != nil {
		log.Println(err)
	}
	defer tempFile.Close()
	data.TempFileName = tempFile.Name()
	data.MediaType = strings.Split(mt, "/")[0]

	tempFile.Write(fileBytes)
}
func init() {
	readDB()
	err := os.Mkdir("./public/temp", 0777)
	if err != nil {
		log.Println(err)
	}
}

func readDB() {
	content, err := os.ReadFile("JSON_DB.json")
	if err != nil {
		// log.Println(err)
	}

	if len(content) > 0 {
		var items []*item
		err := json.Unmarshal(content, &items)
		if err != nil {
			log.Println(err)
		}

		slices.Reverse(items)

		// stream = append(stream, items...)
		stream = items
		for _, item := range stream {
			itemsMap[item.ID] = item
		}
	}
}

func saveJSON() {
	// f, err := os.OpenFile("JSON_DB.json", os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)
	f, err := os.Create("JSON_DB.json")
	if err != nil {
		log.Println(err)
	}

	defer f.Close()

	var stream_ []*item = make([]*item, len(stream))
	copy(stream_, stream)
	slices.Reverse(stream_)
	b, err := json.Marshal(stream_)
	if err != nil {
		log.Println(err)
	}

	if _, err = f.WriteString(string(b)); err != nil {
		log.Println(err)
	}

	readDB()
}

type item struct {
	FileElement  string    `json:"FileElement"`
	Title        string    `json:"Title"`
	Article      string    `json:"Article"`
	ID           string    `json:"ID"`
	TS           time.Time `json:"TS"`
	MediaType    string    `json:"mediaType"`
	TempFileName string    `json:"tempFileName"`
}
