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

func sendMsg(bot *tgbotapi.BotAPI, text string, receiver int64) (tgbotapi.Message, error) {
	msg := tgbotapi.NewMessage(receiver, text)
	return bot.Send(msg)
}

func runBot() {
	if err != nil {
		fmt.Println(err)
	}

	log.Printf("Authorized on account %s\n", bot.Self.UserName)

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
		switch update.Message.Text {
		case "/start":
			if !contains(GetChatIDs(), update.Message.Chat.ID) {
				AddChatID(update.Message.Chat.ID)
			}
			text := "/help: 명령어 목록\n/current: 현재 정보 받기\n\n"
			text += "코로나19 관련 데이터 변동시 메시지를 보내드립니다."
			sendMsg(bot, text, update.Message.Chat.ID)
		case "/start@KOR_corona19_status_robot":
			if !contains(GetChatIDs(), update.Message.Chat.ID) {
				AddChatID(update.Message.Chat.ID)
			}
			text := "/help: 명령어 목록\n/current: 현재 정보 받기\n\n"
			text += "코로나19 관련 데이터 변동시 메시지를 보내드립니다."
			sendMsg(bot, text, update.Message.Chat.ID)
		case "/help":
			text := "/help: 명령어 목록\n/current: 현재 정보 받기\n\n"
			text += "소스코드: https://github.com/code-yeongyu/corona19-telegram-notifier-bot"
			sendMsg(bot, text, update.Message.Chat.ID)
		case "/help@KOR_corona19_status_robot":
			text := "/help: 명령어 목록\n/current: 현재 정보 받기\n\n"
			text += "소스코드: https://github.com/code-yeongyu/corona19-telegram-notifier-bot"
			sendMsg(bot, text, update.Message.Chat.ID)
		case "/current":
			numbers := GetRecentNumbers()
			text := fmt.Sprintf("확진자: %d명\n사망자: %d\n완치자: %d명\n\n현재 감염자: %d명", numbers["confirmed"], numbers["death"], numbers["cured"], numbers["confirmed"]-numbers["death"]-numbers["cured"])

			sendMsg(bot, text, update.Message.Chat.ID)
		case "/current@KOR_corona19_status_robot":
			numbers := GetRecentNumbers()
			text := fmt.Sprintf("확진자: %d명\n사망자: %d\n완치자: %d명\n\n현재 감염자: %d명", numbers["confirmed"], numbers["death"], numbers["cured"], numbers["confirmed"]-numbers["death"]-numbers["cured"])

			sendMsg(bot, text, update.Message.Chat.ID)
		}
	}
}

func alertIfDiff() error {
	recent := GetRecentNumbers()
	numbers := GetNumbersFromNaver()

	recentTotal := 0
	numbersTotal := 0

	//validate numbers and recent
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
	t := time.Now()
	if recentTotal == numbersTotal {
		return nil
	}
	fmt.Printf("%s ", t.Format("02-Jan 15:04:05"))
	fmt.Println(numbers) // for monitoring

	AddNumbers(numbers)

	// put values
	confirmed := numbers["confirmed"]
	death := numbers["death"]
	cured := numbers["cured"]

	text := fmt.Sprintf("확진자: %d명\n사망자: %d\n완치자: %d명\n\n현재 감염자: %d명\n\n확진자 증가수: %d명\n사망자 증가수: %d명\n완치자 증가수: %d명",
		confirmed, death, cured, confirmed-death-cured, confirmed-recent["confirmed"], death-recent["death"], cured-recent["cured"])
	chatIDs := GetChatIDs()

	startTime := time.Now()

	for i := range chatIDs {
		_, err := sendMsg(bot, text, chatIDs[i])
		if err != nil {
			switch err.Error() {
			case "Bad Request: chat not found":
				RemoveChatID(chatIDs[i])
			default:
				time.Sleep(4 * time.Second)
			}
		}
	}
	elapsedTime := time.Since(startTime)
	fmt.Printf("%f took to send to %d people.\n", elapsedTime.Seconds(), len(chatIDs))
	return nil
}

func main() {
	go runBot()
	for {
		err := alertIfDiff()
		if err == nil {
			time.Sleep(5 * time.Minute)
		} else {
			fmt.Println(err)
		}
	}
}
