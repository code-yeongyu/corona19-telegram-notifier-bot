package main

import (
	"fmt"
	"log"
	"os"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

func contains(arr []int64, num int64) bool {
	for _, v := range arr {
		if v == num {
			return true
		}
	}
	return false
}

var bot, err = tgbotapi.NewBotAPI(os.Getenv("TELEGRAM_CODE"))

func runBot() {
	if err != nil {
		fmt.Println(err)
	}

	log.Printf("Authorized on account %s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates, err := bot.GetUpdatesChan(u)
	if err != nil {
		return
	}

	for update := range updates {
		if update.Message == nil { // ignore any non-Message Updates
			continue
		}

		log.Printf("[%s] %s\n", update.Message.From.UserName, update.Message.Text)

		if update.Message.Text == "/start" {
			if contains(GetChatIDs(), update.Message.Chat.ID) {
				continue
			}

			AddChatID(update.Message.Chat.ID)
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "5분마다 http://ncov.mohw.go.kr/index_main.jsp 를 확인하여 데이터에 변동이 있으면 메시지를 보내드립니다.")
			bot.Send(msg)
		} else if update.Message.Text == "/current" {
			numbers := GetRecentNumbers()
			text := fmt.Sprintf("현재 총 %d명의 확진자가 발생하였고, %d명이 사망하셨으며, %d명이 완치되어 %d명이 현재 감염자입니다.", numbers["confirmed"], numbers["death"], numbers["cured"], numbers["confirmed"]-numbers["death"]-numbers["cured"])
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, text)
			bot.Send(msg)
		}
	}
}

func alertIfDiff() {
	recent := GetRecentNumbers()
	numbers := GetNumbers()

	recentTotal := 0
	numbersTotal := 0

	for _, num := range recent {
		recentTotal += num
	}
	for _, num := range numbers {
		numbersTotal += num
	}
	if recentTotal == numbersTotal {
		return
	}
	AddNumbers(numbers)
	confirmed := numbers["confirmed"]
	death := numbers["death"]
	cured := numbers["cured"]

	text := fmt.Sprintf("현재 총 %d명의 확진자가 발생하였고, %d명이 사망하셨으며, %d명이 완치되어 %d명이 현재 감염자입니다.\n\n기존의 데이터에 비해 %d명의 확진자가 추가 되었고, %d명의 사망자가 추가되었으며, %d명의 완치자가 추가되었습니다.",
		confirmed, death, cured, confirmed-death-cured, confirmed-recent["confirmed"], death-recent["death"], cured-recent["cured"])
	chatIDs := GetChatIDs()
	fmt.Println(text)
	for _, chatID := range chatIDs {
		msg := tgbotapi.NewMessage(chatID, text)
		_, err = bot.Send(msg)
		if err != nil {
			RemoveChatID(chatID)
		}
	}
}

func main() {
	go runBot()
	for {
		alertIfDiff()
		time.Sleep(5 * time.Minute)
	}
}