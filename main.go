package main

import (
	"log"
	"time"

	"github.com/tucnak/telebot"
)

const (
	// TOKEN is telegram bot token
	TOKEN = ""
)

var admins = []int{}

func main() {
	bot, err := telebot.NewBot(TOKEN)
	if err != nil {
		log.Fatal(err)
	}

	bot.Messages = make(chan telebot.Message)

	endSignal := make(chan bool)

	go func() {
		for message := range bot.Messages {
			log.Println(message.Sender, message.Text)

			user := GetUser(message.Sender)
			if !user.Exists() {
				user.User = message.Sender
				user.Status = USER_DEFAULT
				user.Create()
			}

			text := message.Text
			if text != "/cancel" && text != "ارسال پیام" && text != "ارسال رزومه" && user.Status == USER_DEFAULT {
				bot.SendMessage(message.Sender, "برای ثبت دیدگاه خود لطفا روی گزینه 'ارسال پیام' کلیک نمایید"+"\n"+"برای ارسال رزومه خود به کمیته فنی ،‌ روی گزینه 'ارسال رزومه' کلیک کنید", &telebot.SendOptions{
					ReplyMarkup: telebot.ReplyMarkup{
						CustomKeyboard: [][]string{
							{
								"ارسال پیام", "ارسال رزومه",
							},
						},
						ResizeKeyboard: true,
					},
				})
			} else if text == "ارسال پیام" { // ارسال پیام
				user.Status = USER_REPORT
				user.Save()
				bot.SendMessage(message.Sender, "متن خود را بنویسید سپس ارسال کنید ، کمیته فنی پیام شما را خواهد خواند"+"\n\n"+"در صورت انصراف روی گزینه /cancel کلیک کنید", &telebot.SendOptions{
					ReplyMarkup: telebot.ReplyMarkup{
						HideCustomKeyboard: true,
					},
				})
			} else if text == "ارسال رزومه" { // ارسال رزومه
				user.Status = USER_RESUME
				user.Save()
				bot.SendMessage(message.Sender, "رزومه خود را بنویسید سپس ارسال کنید ، کمیته فنی آن را بررسی خواهد کرد"+"\n\n"+"در صورت انصراف روی گزینه /cancel کلیک کنید", &telebot.SendOptions{
					ReplyMarkup: telebot.ReplyMarkup{
						HideCustomKeyboard: true,
					},
				})
			} else if text == "/cancel" {
				bot.SendMessage(message.Sender, "برای ثبت دیدگاه خود لطفا روی گزینه 'ارسال پیام' کلیک نمایید"+"\n"+"برای ارسال رزومه خود به کمیته فنی ،‌ روی گزینه 'ارسال رزومه' کلیک کنید", &telebot.SendOptions{
					ReplyMarkup: telebot.ReplyMarkup{
						CustomKeyboard: [][]string{
							{
								"ارسال پیام", "ارسال رزومه",
							},
						},
						ResizeKeyboard: true,
					},
				})
				user.Status = USER_DEFAULT
				user.Save()
			} else {
				for _, admin := range admins {
					userAdmin := telebot.User{ID: admin}
					if user.Status == USER_RESUME {
						bot.SendMessage(userAdmin, "ارسال رزومه", nil)
					} else if user.Status == USER_REPORT {
						bot.SendMessage(userAdmin, "ارسال پیام", nil)
					}
					bot.ForwardMessage(userAdmin, message)
				}

				bot.SendMessage(message.Sender, "با تشکر از شما ، پیام شما مورد بررسی قرار خواهد گرفت "+"\n\n"+"برای ثبت دیدگاه خود لطفا روی گزینه 'ارسال پیام' کلیک نمایید"+"\n"+"برای ارسال رزومه خود به کمیته فنی ،‌ روی گزینه 'ارسال رزومه' کلیک کنید", &telebot.SendOptions{
					ReplyMarkup: telebot.ReplyMarkup{
						CustomKeyboard: [][]string{
							{
								"ارسال پیام", "ارسال رزومه",
							},
						},
						ResizeKeyboard: true,
					},
				})

				user.Status = USER_DEFAULT
				user.Save()
			}
		}
		endSignal <- true
	}()

	bot.Start(time.Second)
	<-endSignal
}
