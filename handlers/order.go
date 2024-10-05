package handlers

// import (
// 	"bot/models"
// 	"bot/repository"

// 	// tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
// )

// func HandleNewOrder(bot *tgbotapi.BotAPI, update tgbotapi.Update) {
// 	userID := update.Message.From.ID
// 	order := models.Order{
// 		UserID:     int(userID), // Convert userID from int64 to int
// 		TotalPrice: 0.0,         // Set the total price or calculate based on selected items
// 		Status:     "Pending",
// 	}

// 	err := repository.SaveOrder(order)
// 	if err != nil {
// 		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Failed to place the order. Please try again.")
// 		bot.Send(msg)
// 		return
// 	}

// 	msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Your order has been placed!")
// 	bot.Send(msg)
// }

// func HandleOrderCommand(bot *tgbotapi.BotAPI, chatID int64) {
// 	msg := tgbotapi.NewMessage(chatID, "Please select your items from the menu.")
// 	bot.Send(msg)
// }
