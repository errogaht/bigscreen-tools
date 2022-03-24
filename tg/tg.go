package tg

//tg means telegram, all stuff related to telegram like bot api, http server for bot etc.

/*
func telApiPostReq(action string) {
	telegramToken := getTelegramToken()
	url := "https://api.telegram.org/bot" + telegramToken + "/" + action

	reqBody := fmt.Sprintf(`{"started": "%s", "timeSpentSeconds": %d}`, date+"T00:00:00.000+0000", secondsSpent)

	req, err := http.NewRequest("POST", url, strings.NewReader(reqBody))
	req.Header.Add("Authorization", "Bearer "+bearerToken)
	//req.Header.Add("content-type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err.Error())
	}

	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			panic(err.Error())
		}
	}(resp.Body)

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err.Error())
	}

	if resp.StatusCode != 200 && resp.StatusCode != 201 {
		panic(fmt.Sprintf("status code: %d, body: %s", resp.StatusCode, body))
	}

}*/
func getUpdates() {
	/*telegramToken := getTelegramToken()

	bot, err := tgbotapi.NewBotAPI(telegramToken)
	if err != nil {
		log.Panic(err)
	}

	bot.Debug = true

	log.Printf("Authorized on account %s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := bot.GetUpdatesChan(u)

	for update := range updates {
		if update.Message != nil { // If we got a message
			log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)

			msg := tgbotapi.NewMessage(update.Message.Chat.ID, update.Message.Text)
			msg.ReplyToMessageID = update.Message.MessageID

			bot.Send(msg)
		}
	}*/
}

//func setWebhook() {
//	bot, err := tgbotapi.NewBotAPI(telegramToken)
//	if err != nil {
//		log.Panic(err)
//	}
//
//	params := make(tgbotapi.Params)
//	params["url"] = "https://34e0-85-174-195-48.ngrok.io"
//
//	request, err := bot.MakeRequest("setWebhook", params)
//	if err != nil {
//		return
//	}
//	_ = len(request.Description)
//}
