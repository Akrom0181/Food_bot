package handlers

// import (
// 	"bot/models"
// 	"bot/repository"
// 	"net/http"

// 	"github.com/gin-gonic/gin"
// 	// tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
// )

// func HandleStart(bot *tgbotapi.BotAPI, update tgbotapi.Update) {
// 	msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Welcome to the Food Delivery Bot!")
// 	bot.Send(msg)
// }

// func HandleUnknown(bot *tgbotapi.BotAPI, update tgbotapi.Update) {
// 	msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Unknown command. Please use /menu to see the options.")
// 	bot.Send(msg)
// }

// func HandleUserCommand(bot *tgbotapi.BotAPI, chatID int64) {
// 	msg := tgbotapi.NewMessage(chatID, "Welcome! Use /menu to see our offerings.")
// 	bot.Send(msg)
// }

// type UserHandler struct {
// 	UserRepo *repository.UserRepository
// }

// // NewUserHandler creates a new instance of UserHandler
// func NewUserHandler(userRepo *repository.UserRepository) *UserHandler {
// 	return &UserHandler{UserRepo: userRepo}
// }

// // CreateUser handles user registration
// func (h *UserHandler) CreateUser(c *gin.Context) {
// 	var user models.User
// 	if err := c.ShouldBindJSON(&user); err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
// 		return
// 	}

// 	// Save the user
// 	if err := h.UserRepo.SaveUser(user); err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
// 		return
// 	}

// 	c.JSON(http.StatusCreated, user)
// }
