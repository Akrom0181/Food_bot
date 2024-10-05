package handlers

// import (
// 	"bot/repository"
// 	"fmt"

// 	// tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
// )

// func HandleShowMenu(bot *tgbotapi.BotAPI, update tgbotapi.Update) {
// 	menuItems, err := repository.GetMenuItems()
// 	if err != nil {
// 		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Failed to retrieve menu.")
// 		bot.Send(msg)
// 		return
// 	}

// 	menuMessage := "Menu:\n"
// 	for _, item := range menuItems {
// 		menuMessage += item.Name + " - $" + fmt.Sprintf("%.2f", item.Price) + "\n"
// 	}

// 	msg := tgbotapi.NewMessage(update.Message.Chat.ID, menuMessage)
// 	bot.Send(msg)
// }

// func HandleMenuCommand(bot *tgbotapi.BotAPI, chatID int64) {
// 	// Replace this with your actual menu fetching logic
// 	menu := "1. Pizza\n2. Burger\n3. Salad\n\nUse /order to place an order."
// 	msg := tgbotapi.NewMessage(chatID, menu)
// 	bot.Send(msg)
// }
