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

		log.Printf("[%s](%d) %s\n", update.Message.From.UserName, update.Message.Chat.ID, update.Message.Text)

		if update.Message.Text == "/start" {
			if !contains(GetChatIDs(), update.Message.Chat.ID) {
				AddChatID(update.Message.Chat.ID)
			}

			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "5분마다 확진자의 데이터를 확인하여 데이터에 변동이 있으면 메시지를 보내드립니다.\n소스코드는 다음의 링크에서 확인 하실 수 있습니다: https://github.com/code-yeongyu/corona19-telegram-notifier-bot")
			bot.Send(msg)
		} else if update.Message.Text == "/current" {
			numbers := GetRecentNumbers()
			text := fmt.Sprintf("확진자: %d명\n사망자: %d\n완치자: %d명\n\n따라서 현재 감염자: %d명", numbers["confirmed"], numbers["death"], numbers["cured"], numbers["confirmed"]-numbers["death"]-numbers["cured"])
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, text)
			bot.Send(msg)
		} else {
			text := "/current: 현재 정보 받기\n\n"
			text += "5분마다 확진자의 데이터를 확인하여 데이터에 변동이 있으면 메시지를 보내드립니다.\n소스코드는 다음의 링크에서 확인 하실 수 있습니다: https://github.com/code-yeongyu/corona19-telegram-notifier-bot"
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, text)
			bot.Send(msg)
		}
	}
}

func alertIfDiff() error {
	recent := GetRecentNumbers()
	numbers := GetNumbers()

	recentTotal := 0
	numbersTotal := 0

	if len(numbers) == 0 || len(recent) == 0 {
		return fmt.Errorf("something includes 0")
	}

	for _, num := range recent {
		if num == 0 {
			return fmt.Errorf("recent includes 0")
		}
		recentTotal += num
	}
	for _, num := range numbers {
		if num == 0 {
			return fmt.Errorf("numbers includes 0")
		}
		numbersTotal += num
	}
	if recentTotal == numbersTotal {
		return nil
	}
	AddNumbers(numbers)
	confirmed := numbers["confirmed"]
	death := numbers["death"]
	cured := numbers["cured"]

	text := fmt.Sprintf("확진자: %d명\n사망자: %d\n완치자: %d명\n\n따라서 현재 감염자: %d명\n\n확진자 증가수: %d명\n사망자 증가수: %d명\n완치자 증가수: %d명",
		confirmed, death, cured, confirmed-death-cured, confirmed-recent["confirmed"], death-recent["death"], cured-recent["cured"])
	chatIDs := GetChatIDs()
	fmt.Println(numbers)
	for _, chatID := range chatIDs {
		msg := tgbotapi.NewMessage(chatID, text)
		_, err = bot.Send(msg)
		if err != nil {
			RemoveChatID(chatID)
		}
	}
	return nil
}

func main() {
	go runBot()
	for {
		err := alertIfDiff()
		if err == nil {
			time.Sleep(5 * time.Minute)
		}
		fmt.Println("fuck")
	}
}
