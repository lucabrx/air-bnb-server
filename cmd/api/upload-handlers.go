package main

import (
	"net/http"
	"strings"

	"github.com/google/uuid"
	"github.com/minio/minio-go/v7"
)

func (app *application) uploadImageHandler(w http.ResponseWriter, r *http.Request) {
	err := r.ParseMultipartForm(10 << 20)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	file, handler, err := r.FormFile("file")
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}
	defer file.Close()

	c, err := app.aws.Client()
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	id, err := uuid.NewUUID()
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
	ext := strings.Split(handler.Header.Get("Content-Type"), "/")[1]
	name := id.String() + "." + ext
	url := "https://air-bnb-clone-luka.s3.eu-central-1.amazonaws.com/" + name

	_, err = c.PutObject(r.Context(), "air-bnb-clone-luka", name, file, handler.Size, minio.PutObjectOptions{ContentType: handler.Header.Get("Content-Type")})
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
	err = app.writeJSON(w, http.StatusCreated, envelope{"url": url}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
}
