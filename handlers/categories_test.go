package handlers

import (
	"github.com/woodhds/vk.service/predictor"
	"testing"
)

func TestSave(t *testing.T) {
	if client, e := predictor.NewClient("http://vk-predict.herokuapp.com:80"); e != nil {
		t.Error(e)
	} else {
		if e := client.SaveMessage(1, 2, "", ""); e != nil {
			t.Error(e)
		}
	}
}

func TestGet(t *testing.T) {
	if client, e := predictor.NewClient("http://vk-predict.herokuapp.com:80"); e != nil {
		t.Error(e)
	} else {
		if resp, e := client.Predict([]*predictor.PredictMessage{
			{OwnerId: -164852303, Id: 2309, Category: "", Text: `🧚‍♀️- Друзья...🤗
			👨‍🔧- Ну наконец-то мы домастерили вот такой вот симпатичный платяной шкафчик. 
			 Давайте же уже его разыграем 🎁
			
			❗- Что нужно сделать👇
			  - Подписываемся на нашу группу
			  - Ставим лайки
			  - Делаем репост
			....Победителя объявим 1-го июня.Всем удачи и хорошего настроения.`},
		}); e != nil {
			t.Error(e)
		} else {
			if len(resp) == 0 || resp[0].Category == "" {
				t.Error("Response empty")
			}
		}
	}
}
