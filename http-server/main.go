package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	api "xtz-Server-Assignment-TikTokTechImmersion/api"
	db "xtz-Server-Assignment-TikTokTechImmersion/http-server/db"
)

type Server struct {
	Database *db.Database
}

func main() {
	database, err := db.NewDatabase()
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer database.DB.Close()

	s := &Server{
		Database: database,
	}

	http.HandleFunc("/send", func(w http.ResponseWriter, r *http.Request) {
		r.ParseForm()

		req := &api.SendRequest{
			Chat:   r.Form.Get("chat"),
			Text:   r.Form.Get("text"),
			Sender: r.Form.Get("sender"),
		}

		_, err := s.Database.Send(context.Background(), req)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	})

	http.HandleFunc("/pull", func(w http.ResponseWriter, r *http.Request) {
		r.ParseForm()

		limit, _ := strconv.ParseInt(r.Form.Get("limit"), 10, 32)
		cursor, _ := strconv.ParseInt(r.Form.Get("cursor"), 10, 64)

		req := &api.PullRequest{
			Chat:   r.Form.Get("chat"),
			Limit:  int32(limit),
			Cursor: cursor,
		}

		resp, err := s.Database.Pull(context.Background(), req)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		json.NewEncoder(w).Encode(resp.Messages)
	})

	log.Println("Starting HTTP server on port 8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
