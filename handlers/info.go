package handlers

import (
	"github.com/woodhds/vk.service/database"
	"net/http"
	"strconv"
)

func InfoHandler(messageQueryService database.MessagesQueryService) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		count := messageQueryService.GetTotalCount()
		writer.Write([]byte(strconv.Itoa(count)))
		writer.WriteHeader(http.StatusOK)
	})
}