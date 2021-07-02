package handlers

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
	"net/url"
	"strconv"
)

type saveRequest struct {
	Category string `json:"category"`
}

type MessageSaver interface {
	SaveMessage(owner int, id int, text string, category string) error
}

type MessageSaveClient struct {
	httpClient *http.Client
	host       string
}

func (m MessageSaveClient) SaveMessage(owner int, id int, text string, category string) error {
	reqData := make(map[string]interface{})
	reqData["ownerId"] = owner
	reqData["id"] = id
	reqData["text"] = text
	reqData["category"] = category
	b, _ := json.Marshal(reqData)
	u, _ := url.Parse(m.host)
	u.Path = "/predict"
	if req, e := http.NewRequest(http.MethodPut, u.String(), bytes.NewBuffer(b)); e == nil {
		req.Header.Add("content-type", "application/json")
		if resp, e := m.httpClient.Do(req); e != nil {
			fmt.Println(resp)
			return e
		} else {
			if resp.StatusCode != http.StatusOK {
				return errors.New(fmt.Sprintf("Server responded with status code %d", resp.StatusCode))
			}
		}
	} else {
		return e
	}

	return nil
}

func MessageSaveHandler(host string, conn *sql.DB) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		var owner int
		var messageId int
		if ownerId, e := strconv.Atoi(vars["ownerId"]); e == nil {
			owner = ownerId
		}

		if id, e := strconv.Atoi(vars["id"]); e == nil {
			messageId = id
		}

		var data saveRequest
		json.NewDecoder(r.Body).Decode(&data)

		if owner == 0 || messageId == 0 {
			rw.WriteHeader(http.StatusBadRequest)
			return
		}

		var text string

		if d, e:=conn.Query("SELECT Text from messages where OwnerId = $1 and Id = $2", owner, messageId); e != nil {
			rw.WriteHeader(http.StatusBadRequest)
			return
		} else {
			for d.Next() {
				d.Scan(&text)
			}
		}

		if client, e := NewClient(host); e != nil {
			fmt.Println(e)
			rw.WriteHeader(http.StatusBadRequest)
		} else {
			if e := client.SaveMessage(owner, messageId, text, data.Category); e != nil {
				rw.WriteHeader(http.StatusBadRequest)
			}
			Json(rw, true)
		}

	})
}

func NewClient(host string) (MessageSaver, error) {
	return &MessageSaveClient{
		httpClient: &http.Client{},
		host: host,
	}, nil
}
