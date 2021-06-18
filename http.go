package main

import (
	"encoding/json"
	"net/http"
)

func Json(rw http.ResponseWriter, data interface{}) {
	if j, e := json.Marshal(data); e != nil {
		rw.WriteHeader(http.StatusBadRequest)
	} else {
		rw.Header().Add("Content-type", "application/json")
		rw.Write(j)
		rw.WriteHeader(http.StatusOK)
	}
}
