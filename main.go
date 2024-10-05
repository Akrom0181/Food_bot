// package main

// import (
// 	"bot/config"
// 	"bot/handlers"
// 	"database/sql"
// 	"fmt"
// 	"log"

// 	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
// 	"github.com/joho/godotenv"

// 	_ "github.com/lib/pq" // Import the PostgreSQL driver
// )

// var db *sql.DB

// func main() {
// 	err := godotenv.Load()
// 	if err != nil {
// 		log.Fatal("Error loading .env file")
// 	}

// 	bot, err := tgbotapi.NewBotAPI("7760253527:AAGNyUVk3NZ1f2RCpdfy_2rNuUaxYKQX2-4")
// 	if err != nil {
// 		log.Panic(err)
// 	}

// 	bot.Debug = true
// 	log.Printf("Authorized on account %s", bot.Self.UserName)

// 	// Initialize database connection
// 	initDB()

// 	updates := bot.GetUpdatesChan(tgbotapi.NewUpdate(0))

// 	for update := range updates {
// 		if update.Message != nil {
// 			handleMessage(bot, update.Message)
// 		} else if update.CallbackQuery != nil {
// 			handleCallbackQuery(bot, update.CallbackQuery)
// 		}
// 	}
// }

// func initDB() {
// 	cfg := config.LoadConfig()
// 	connectionString := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
// 		cfg.DBHost, cfg.DBPort, cfg.DBUser, cfg.DBPassword, cfg.DBName)

// 	var err error
// 	db, err = sql.Open("postgres", connectionString)
// 	if err != nil {
// 		log.Fatal("Failed to connect to database:", err)
// 	}

// 	err = db.Ping()
// 	if err != nil {
// 		log.Fatal("Failed to ping database:", err)
// 	}

// 	log.Println("Database connected successfully!")
// }

// func handleMessage(bot *tgbotapi.BotAPI, message *tgbotapi.Message) {
// 	switch message.Command() {
// 	case "start":
// 		handlers.HandleUserCommand(bot, message.Chat.ID)
// 	case "menu":
// 		handlers.HandleMenuCommand(bot, message.Chat.ID)
// 	case "order":
// 		handlers.HandleOrderCommand(bot, message.Chat.ID)
// 	default:
// 		msg := tgbotapi.NewMessage(message.Chat.ID, "I don't understand that command. Please use /menu or /order.")
// 		bot.Send(msg)
// 	}
// }

// func handleCallbackQuery(bot *tgbotapi.BotAPI, callback *tgbotapi.CallbackQuery) {
// 	msg := tgbotapi.NewMessage(callback.Message.Chat.ID, "You clicked: "+callback.Data)
// 	bot.Send(msg)

//		// Acknowledge the callback query
//		_, err := bot.Request(tgbotapi.CallbackConfig{
//			CallbackQueryID: callback.ID,
//			Text:            "Processing...",
//			ShowAlert:       false,
//		})
//		if err != nil {
//			log.Println("Error answering callback query:", err)
//		}
//	}
// package main

// import (
// 	"database/sql"
// 	"fmt"
// 	"log"
// 	"strconv"
// 	"strings"

// 	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
// 	_ "github.com/lib/pq"
// )

// // Data structures
// type MenuItem struct {
// 	ID    int
// 	Name  string
// 	Price float64
// }

// type OrderItem struct {
// 	Item     MenuItem
// 	Quantity int
// }

// type Order struct {
// 	ID     int
// 	UserID int64
// 	Items  []OrderItem
// 	Total  float64
// }

// // Global variables
// var (
// 	db  *sql.DB
// 	bot *tgbotapi.BotAPI
// )

// // "shahzod",
// // "food_ordering_bot",
// // "1",
// // "localhost",
// // "5432",
// // Database functions
// func initDB() {
// 	var err error
// 	// Fetch database credentials from environment variables
// 	dbUser := "shahzod"
// 	dbName := "food_ordering_bot"
// 	dbPassword := "1"
// 	dbSSLMode := "disable"

// 	if dbUser == "" || dbName == "" || dbPassword == "" || dbSSLMode == "" {
// 		log.Fatal("Database environment variables (DB_USER, DB_NAME, DB_PASSWORD, DB_SSLMODE) are not set")
// 	}

// 	connStr := fmt.Sprintf("user=%s dbname=%s sslmode=%s password=%s",
// 		dbUser, dbName, dbSSLMode, dbPassword)
// 	db, err = sql.Open("postgres", connStr)
// 	if err != nil {
// 		log.Fatalf("Failed to open database: %v", err)
// 	}

// 	err = db.Ping()
// 	if err != nil {
// 		log.Fatalf("Failed to ping database: %v", err)
// 	}

// 	createTables()
// 	insertInitialData()
// }

// func createTables() {
// 	queries := `
// 	CREATE TABLE IF NOT EXISTS menu_items (
// 		id SERIAL PRIMARY KEY,
// 		name TEXT NOT NULL UNIQUE,
// 		price REAL NOT NULL
// 	);
// 	CREATE TABLE IF NOT EXISTS orders (
// 		id SERIAL PRIMARY KEY,
// 		user_id BIGINT NOT NULL,
// 		total REAL NOT NULL
// 	);
// 	CREATE TABLE IF NOT EXISTS order_items (
// 		order_id INTEGER REFERENCES orders(id) ON DELETE CASCADE,
// 		item_id INTEGER REFERENCES menu_items(id),
// 		quantity INTEGER NOT NULL,
// 		PRIMARY KEY (order_id, item_id)
// 	);
// 	CREATE TABLE IF NOT EXISTS carts (
// 		user_id BIGINT PRIMARY KEY,
// 		items TEXT NOT NULL
// 	);
// 	`

// 	_, err := db.Exec(queries)
// 	if err != nil {
// 		log.Fatalf("Failed to create tables: %v", err)
// 	}
// }

// func insertInitialData() {
// 	// Insert menu items only if they don't already exist
// 	_, err := db.Exec(`
// 		INSERT INTO menu_items (name, price) VALUES
// 		('Osh', 25000),
// 		('Lagmon', 22000),
// 		('Shashlik', 15000),
// 		('Cola', 5000)
// 		ON CONFLICT (name) DO NOTHING;
// 	`)
// 	if err != nil {
// 		log.Printf("Error inserting initial data: %v", err)
// 	}
// }

// func getMenuItems() ([]MenuItem, error) {
// 	rows, err := db.Query("SELECT id, name, price FROM menu_items")
// 	if err != nil {
// 		return nil, fmt.Errorf("error querying menu items: %w", err)
// 	}
// 	defer rows.Close()

// 	var items []MenuItem
// 	for rows.Next() {
// 		var item MenuItem
// 		if err := rows.Scan(&item.ID, &item.Name, &item.Price); err != nil {
// 			log.Printf("Error scanning menu item: %v", err)
// 			continue
// 		}
// 		items = append(items, item)
// 	}
// 	return items, nil
// }

// func getOrderItems(orderID int) ([]OrderItem, error) {
// 	rows, err := db.Query(`
// 		SELECT mi.id, mi.name, mi.price, oi.quantity
// 		FROM order_items oi
// 		JOIN menu_items mi ON oi.item_id = mi.id
// 		WHERE oi.order_id = $1
// 	`, orderID)
// 	if err != nil {
// 		return nil, fmt.Errorf("error querying order items: %w", err)
// 	}
// 	defer rows.Close()

// 	var items []OrderItem
// 	for rows.Next() {
// 		var oi OrderItem
// 		if err := rows.Scan(&oi.Item.ID, &oi.Item.Name, &oi.Item.Price, &oi.Quantity); err != nil {
// 			log.Printf("Error scanning order item: %v", err)
// 			continue
// 		}
// 		items = append(items, oi)
// 	}
// 	return items, nil
// }

// func getUserOrders(userID int64) ([]Order, error) {
// 	rows, err := db.Query(`
// 		SELECT id, total FROM orders
// 		WHERE user_id = $1
// 		ORDER BY id DESC
// 	`, userID)
// 	if err != nil {
// 		return nil, fmt.Errorf("error querying orders: %w", err)
// 	}
// 	defer rows.Close()

// 	var orders []Order
// 	for rows.Next() {
// 		var order Order
// 		if err := rows.Scan(&order.ID, &order.Total); err != nil {
// 			log.Printf("Error scanning order: %v", err)
// 			continue
// 		}
// 		order.UserID = userID
// 		items, err := getOrderItems(order.ID)
// 		if err != nil {
// 			log.Printf("Error getting order items: %v", err)
// 			continue
// 		}
// 		order.Items = items
// 		orders = append(orders, order)
// 	}
// 	return orders, nil
// }

// func addToCart(userID int64, itemID int) error {
// 	// Fetch current cart
// 	var itemsStr string
// 	err := db.QueryRow("SELECT items FROM carts WHERE user_id = $1", userID).Scan(&itemsStr)
// 	if err != nil && err != sql.ErrNoRows {
// 		return fmt.Errorf("error fetching cart: %w", err)
// 	}

// 	var itemsMap map[int]int
// 	if itemsStr != "" {
// 		itemsMap = make(map[int]int)
// 		pairs := strings.Split(itemsStr, ",")
// 		for _, pair := range pairs {
// 			parts := strings.Split(pair, ":")
// 			if len(parts) != 2 {
// 				continue
// 			}
// 			id, err := strconv.Atoi(parts[0])
// 			if err != nil {
// 				continue
// 			}
// 			qty, err := strconv.Atoi(parts[1])
// 			if err != nil {
// 				continue
// 			}
// 			itemsMap[id] = qty
// 		}
// 	} else {
// 		itemsMap = make(map[int]int)
// 	}

// 	// Increment quantity
// 	itemsMap[itemID]++

// 	// Reconstruct items string
// 	var newItems []string
// 	for id, qty := range itemsMap {
// 		newItems = append(newItems, fmt.Sprintf("%d:%d", id, qty))
// 	}
// 	newItemsStr := strings.Join(newItems, ",")

// 	// Upsert cart
// 	_, err = db.Exec(`
// 		INSERT INTO carts (user_id, items)
// 		VALUES ($1, $2)
// 		ON CONFLICT (user_id) DO UPDATE SET items = EXCLUDED.items
// 	`, userID, newItemsStr)
// 	if err != nil {
// 		return fmt.Errorf("error updating cart: %w", err)
// 	}

// 	return nil
// }

// func clearCart(userID int64) error {
// 	_, err := db.Exec("DELETE FROM carts WHERE user_id = $1", userID)
// 	return err
// }

// func placeOrder(userID int64) error {
// 	// Fetch cart
// 	var itemsStr string
// 	err := db.QueryRow("SELECT items FROM carts WHERE user_id = $1", userID).Scan(&itemsStr)
// 	if err != nil {
// 		if err == sql.ErrNoRows {
// 			return fmt.Errorf("cart is empty")
// 		}
// 		return fmt.Errorf("error fetching cart: %w", err)
// 	}

// 	// Parse items
// 	itemsMap := make(map[int]int)
// 	pairs := strings.Split(itemsStr, ",")
// 	for _, pair := range pairs {
// 		parts := strings.Split(pair, ":")
// 		if len(parts) != 2 {
// 			continue
// 		}
// 		id, err := strconv.Atoi(parts[0])
// 		if err != nil {
// 			continue
// 		}
// 		qty, err := strconv.Atoi(parts[1])
// 		if err != nil {
// 			continue
// 		}
// 		itemsMap[id] = qty
// 	}

// 	if len(itemsMap) == 0 {
// 		return fmt.Errorf("cart is empty")
// 	}

// 	// Calculate total
// 	var total float64
// 	for id, qty := range itemsMap {
// 		var price float64
// 		err := db.QueryRow("SELECT price FROM menu_items WHERE id = $1", id).Scan(&price)
// 		if err != nil {
// 			return fmt.Errorf("error fetching price for item %d: %w", id, err)
// 		}
// 		total += price * float64(qty)
// 	}

// 	// Insert order
// 	var orderID int
// 	err = db.QueryRow(`
// 		INSERT INTO orders (user_id, total)
// 		VALUES ($1, $2)
// 		RETURNING id
// 	`, userID, total).Scan(&orderID)
// 	if err != nil {
// 		return fmt.Errorf("error inserting order: %w", err)
// 	}

// 	// Insert order items
// 	for id, qty := range itemsMap {
// 		_, err := db.Exec(`
// 			INSERT INTO order_items (order_id, item_id, quantity)
// 			VALUES ($1, $2, $3)
// 		`, orderID, id, qty)
// 		if err != nil {
// 			return fmt.Errorf("error inserting order item: %w", err)
// 		}
// 	}

// 	// Clear cart
// 	err = clearCart(userID)
// 	if err != nil {
// 		return fmt.Errorf("error clearing cart: %w", err)
// 	}

// 	return nil
// }

// // Bot functions
// func initBot() {
// 	var err error
// 	// Fetch bot token from environment variables
// 	botToken := "7961675553:AAGPX7myxcS3EBj-Rz7RlSXdrXwXunFbW0Q"
// 	if botToken == "" {
// 		log.Fatal("TELEGRAM_BOT_API_TOKEN is not set")
// 	}

// 	bot, err = tgbotapi.NewBotAPI(botToken)
// 	if err != nil {
// 		log.Fatalf("Failed to create bot: %v", err)
// 	}

// 	bot.Debug = true
// 	log.Printf("Authorized on account %s", bot.Self.UserName)
// }

// func handleUpdates() {
// 	u := tgbotapi.NewUpdate(0)
// 	u.Timeout = 60

// 	updates := bot.GetUpdatesChan(u)

// 	for update := range updates {
// 		if update.Message != nil {
// 			handleMessage(update.Message)
// 		} else if update.CallbackQuery != nil {
// 			handleCallbackQuery(update.CallbackQuery)
// 		}
// 	}
// }

// func handleMessage(message *tgbotapi.Message) {
// 	if message.Text == "/start" {
// 		sendMainMenu(message.Chat.ID)
// 	} else if strings.HasPrefix(message.Text, "/order") {
// 		// Handle placing an order
// 		err := placeOrder(message.From.ID)
// 		if err != nil {
// 			msg := tgbotapi.NewMessage(message.Chat.ID, fmt.Sprintf("Failed to place order: %v", err))
// 			bot.Send(msg)
// 			return
// 		}
// 		msg := tgbotapi.NewMessage(message.Chat.ID, "Your order has been placed successfully!")
// 		bot.Send(msg)
// 	} else {
// 		msg := tgbotapi.NewMessage(message.Chat.ID, "I didn't understand that command. Please use /start to see options.")
// 		bot.Send(msg)
// 	}
// }

// func handleCallbackQuery(query *tgbotapi.CallbackQuery) {
// 	data := query.Data

// 	switch {
// 	case data == "menu":
// 		sendMenu(query.Message.Chat.ID)
// 	case data == "main_menu":
// 		sendMainMenu(query.Message.Chat.ID)
// 	case data == "my_orders":
// 		sendOrders(query.Message.Chat.ID, query.From.ID)
// 	case strings.HasPrefix(data, "add_to_cart_"):
// 		itemIDStr := strings.TrimPrefix(data, "add_to_cart_")
// 		itemID, err := strconv.Atoi(itemIDStr)
// 		if err != nil {
// 			log.Printf("Invalid item ID: %v", err)
// 			answerCallback(query.ID, "Invalid item.")
// 			return
// 		}

// 		err = addToCart(query.From.ID, itemID)
// 		if err != nil {
// 			log.Printf("Error adding to cart: %v", err)
// 			answerCallback(query.ID, "Failed to add to cart.")
// 			return
// 		}

// 		answerCallback(query.ID, "Added to cart!")
// 	case data == "view_cart":
// 		sendCart(query.Message.Chat.ID, query.From.ID)
// 	case data == "checkout":
// 		err := placeOrder(query.From.ID)
// 		if err != nil {
// 			answerCallback(query.ID, fmt.Sprintf("Failed to place order: %v", err))
// 			return
// 		}
// 		answerCallback(query.ID, "Your order has been placed!")
// 		sendMainMenu(query.Message.Chat.ID)
// 	case data == "clear_cart":
// 		err := clearCart(query.From.ID)
// 		if err != nil {
// 			log.Printf("Error clearing cart: %v", err)
// 			answerCallback(query.ID, "Failed to clear cart.")
// 			return
// 		}
// 		answerCallback(query.ID, "Your cart has been cleared.")
// 		sendMainMenu(query.Message.Chat.ID)
// 	default:
// 		answerCallback(query.ID, "Unknown action.")
// 	}
// }

// func answerCallback(callbackID, text string) {
// 	callback := tgbotapi.CallbackConfig{
// 		CallbackQueryID: callbackID,
// 		Text:            text,
// 		ShowAlert:       false,
// 	}

// 	if _, err := bot.Request(callback); err != nil {
// 		log.Printf("Failed to answer callback query: %v", err)
// 	}
// }

// func sendMainMenu(chatID int64) {
// 	keyboard := tgbotapi.NewInlineKeyboardMarkup(
// 		tgbotapi.NewInlineKeyboardRow(
// 			tgbotapi.NewInlineKeyboardButtonData("📖 Menu", "menu"),
// 			tgbotapi.NewInlineKeyboardButtonData("🛒 My Cart", "view_cart"),
// 		),
// 		tgbotapi.NewInlineKeyboardRow(
// 			tgbotapi.NewInlineKeyboardButtonData("📦 My Orders", "my_orders"),
// 		),
// 	)

// 	msg := tgbotapi.NewMessage(chatID, "Welcome to our restaurant bot! Please choose an option:")
// 	msg.ReplyMarkup = keyboard
// 	if _, err := bot.Send(msg); err != nil {
// 		log.Printf("Failed to send main menu: %v", err)
// 	}
// }

// func sendMenu(chatID int64) {
// 	items, err := getMenuItems()
// 	if err != nil {
// 		log.Printf("Error getting menu items: %v", err)
// 		msg := tgbotapi.NewMessage(chatID, "Failed to retrieve menu.")
// 		bot.Send(msg)
// 		return
// 	}

// 	var keyboard [][]tgbotapi.InlineKeyboardButton
// 	for _, item := range items {
// 		button := tgbotapi.NewInlineKeyboardButtonData(
// 			fmt.Sprintf("%s - %.0f so'm", item.Name, item.Price),
// 			fmt.Sprintf("add_to_cart_%d", item.ID),
// 		)
// 		row := tgbotapi.NewInlineKeyboardRow(button)
// 		keyboard = append(keyboard, row)
// 	}

// 	// Add a back button
// 	backButton := tgbotapi.NewInlineKeyboardButtonData("🔙 Back", "main_menu")
// 	keyboard = append(keyboard, tgbotapi.NewInlineKeyboardRow(backButton))

// 	msg := tgbotapi.NewMessage(chatID, "Our Menu:")
// 	msg.ReplyMarkup = tgbotapi.NewInlineKeyboardMarkup(keyboard...)
// 	if _, err := bot.Send(msg); err != nil {
// 		log.Printf("Failed to send menu: %v", err)
// 	}
// }

// func sendOrders(chatID int64, userID int64) {
// 	orders, err := getUserOrders(userID)
// 	if err != nil {
// 		log.Printf("Error getting user orders: %v", err)
// 		msg := tgbotapi.NewMessage(chatID, "Failed to retrieve your orders.")
// 		bot.Send(msg)
// 		return
// 	}

// 	if len(orders) == 0 {
// 		msg := tgbotapi.NewMessage(chatID, "You have no orders yet.")
// 		bot.Send(msg)
// 		return
// 	}

// 	var response strings.Builder
// 	for _, order := range orders {
// 		response.WriteString(fmt.Sprintf("🆔 Order #%d\n💰 Total: %.0f so'm\n\n", order.ID, order.Total))
// 	}

// 	msg := tgbotapi.NewMessage(chatID, response.String())
// 	bot.Send(msg)
// }

// func sendCart(chatID int64, userID int64) {
// 	// Fetch cart
// 	var itemsStr string
// 	err := db.QueryRow("SELECT items FROM carts WHERE user_id = $1", userID).Scan(&itemsStr)
// 	if err != nil {
// 		if err == sql.ErrNoRows {
// 			msg := tgbotapi.NewMessage(chatID, "Your cart is empty.")
// 			bot.Send(msg)
// 			return
// 		}
// 		log.Printf("Error fetching cart: %v", err)
// 		msg := tgbotapi.NewMessage(chatID, "Failed to retrieve your cart.")
// 		bot.Send(msg)
// 		return
// 	}

// 	// Parse items
// 	itemsMap := make(map[int]int)
// 	pairs := strings.Split(itemsStr, ",")
// 	for _, pair := range pairs {
// 		parts := strings.Split(pair, ":")
// 		if len(parts) != 2 {
// 			continue
// 		}
// 		id, err := strconv.Atoi(parts[0])
// 		if err != nil {
// 			continue
// 		}
// 		qty, err := strconv.Atoi(parts[1])
// 		if err != nil {
// 			continue
// 		}
// 		itemsMap[id] = qty
// 	}

// 	if len(itemsMap) == 0 {
// 		msg := tgbotapi.NewMessage(chatID, "Your cart is empty.")
// 		bot.Send(msg)
// 		return
// 	}

// 	// Fetch item details
// 	var response strings.Builder
// 	var total float64
// 	for id, qty := range itemsMap {
// 		var name string
// 		var price float64
// 		err := db.QueryRow("SELECT name, price FROM menu_items WHERE id = $1", id).Scan(&name, &price)
// 		if err != nil {
// 			log.Printf("Error fetching item %d: %v", id, err)
// 			continue
// 		}
// 		response.WriteString(fmt.Sprintf("• %s x%d - %.0f so'm\n", name, qty, price*float64(qty)))
// 		total += price * float64(qty)
// 	}

// 	response.WriteString(fmt.Sprintf("\n💰 Total: %.0f so'm", total))

// 	// Add buttons for checkout, clear cart, and back
// 	keyboard := tgbotapi.NewInlineKeyboardMarkup(
// 		tgbotapi.NewInlineKeyboardRow(
// 			tgbotapi.NewInlineKeyboardButtonData("✅ Checkout", "checkout"),
// 			tgbotapi.NewInlineKeyboardButtonData("🗑️ Clear Cart", "clear_cart"),
// 		),
// 		tgbotapi.NewInlineKeyboardRow(
// 			tgbotapi.NewInlineKeyboardButtonData("🔙 Back to Menu", "main_menu"),
// 		),
// 	)

// 	msg := tgbotapi.NewMessage(chatID, response.String())
// 	msg.ReplyMarkup = keyboard
// 	if _, err := bot.Send(msg); err != nil {
// 		log.Printf("Failed to send cart: %v", err)
// 	}
// }

// // Main function
// func main() {
// 	// Initialize logging
// 	log.SetFlags(log.LstdFlags | log.Lshortfile)

// 	// Initialize database
// 	initDB()
// 	defer db.Close()

// 	// Initialize bot
// 	initBot()

//		// Start handling updates
//		handleUpdates()
//	}
package main

import (
	"database/sql"
	"fmt"
	"log"
	"strconv"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	_ "github.com/lib/pq"
)

// Data structures
type MenuItem struct {
	ID    int
	Name  string
	Price float64
}

type OrderItem struct {
	Item     MenuItem
	Quantity int
}

type Order struct {
	ID     int
	UserID int64
	Items  []OrderItem
	Total  float64
}

type User struct {
	UserID   int64
	Phone    string
	Location string
}

// Global variables
var (
	db  *sql.DB
	bot *tgbotapi.BotAPI
)

// Database functions
func initDB() {
	var err error
	// Fetch database credentials from environment variables
	dbUser := "shahzod"
	dbName := "food_gehc"
	dbPassword := "YuZERp29JK67MAtUJFN83h1e2FLe7TpE"
	dbSSLMode := "dpg-crg01paj1k6c7399vc90-a.ohio-postgres.render.com"

	if dbUser == "" || dbName == "" || dbPassword == "" || dbSSLMode == "" {
		log.Fatal("Database environment variables (DB_USER, DB_NAME, DB_PASSWORD, DB_SSLMODE) are not set")
	}

	connStr := fmt.Sprintf("user=%s dbname=%s sslmode=%s password=%s",
		dbUser, dbName, dbSSLMode, dbPassword)
	db, err = sql.Open("postgres", connStr)
	if err != nil {
		log.Fatalf("Failed to open database: %v", err)
	}

	err = db.Ping()
	if err != nil {
		log.Fatalf("Failed to ping database: %v", err)
	}

	createTables()
	insertInitialData()
}

func createTables() {
	queries := `
	CREATE TABLE IF NOT EXISTS menu_items (
		id SERIAL PRIMARY KEY,
		name TEXT NOT NULL UNIQUE,
		price REAL NOT NULL
	);
	CREATE TABLE IF NOT EXISTS orders (
		id SERIAL PRIMARY KEY,
		user_id BIGINT NOT NULL,
		total REAL NOT NULL
	);
	CREATE TABLE IF NOT EXISTS order_items (
		order_id INTEGER REFERENCES orders(id) ON DELETE CASCADE,
		item_id INTEGER REFERENCES menu_items(id),
		quantity INTEGER NOT NULL,
		PRIMARY KEY (order_id, item_id)
	);
	CREATE TABLE IF NOT EXISTS carts (
		user_id BIGINT PRIMARY KEY,
		items TEXT NOT NULL
	);
	CREATE TABLE IF NOT EXISTS users (
		user_id BIGINT PRIMARY KEY,
		phone TEXT,
		location TEXT
	);
	`

	_, err := db.Exec(queries)
	if err != nil {
		log.Fatalf("Failed to create tables: %v", err)
	}
}

func insertInitialData() {
	// Insert menu items only if they don't already exist
	_, err := db.Exec(`
		INSERT INTO menu_items (name, price) VALUES
		('Osh', 25000),
		('Lagmon', 22000),
		('Shashlik', 15000),
		('Cola', 5000)
		ON CONFLICT (name) DO NOTHING;
	`)
	if err != nil {
		log.Printf("Error inserting initial data: %v", err)
	}
}

func getMenuItems() ([]MenuItem, error) {
	rows, err := db.Query("SELECT id, name, price FROM menu_items")
	if err != nil {
		return nil, fmt.Errorf("error querying menu items: %w", err)
	}
	defer rows.Close()

	var items []MenuItem
	for rows.Next() {
		var item MenuItem
		if err := rows.Scan(&item.ID, &item.Name, &item.Price); err != nil {
			log.Printf("Error scanning menu item: %v", err)
			continue
		}
		items = append(items, item)
	}
	return items, nil
}

func getOrderItems(orderID int) ([]OrderItem, error) {
	rows, err := db.Query(`
		SELECT mi.id, mi.name, mi.price, oi.quantity
		FROM order_items oi
		JOIN menu_items mi ON oi.item_id = mi.id
		WHERE oi.order_id = $1
	`, orderID)
	if err != nil {
		return nil, fmt.Errorf("error querying order items: %w", err)
	}
	defer rows.Close()

	var items []OrderItem
	for rows.Next() {
		var oi OrderItem
		if err := rows.Scan(&oi.Item.ID, &oi.Item.Name, &oi.Item.Price, &oi.Quantity); err != nil {
			log.Printf("Error scanning order item: %v", err)
			continue
		}
		items = append(items, oi)
	}
	return items, nil
}

func getUserOrders(userID int64) ([]Order, error) {
	rows, err := db.Query(`
		SELECT id, total FROM orders
		WHERE user_id = $1
		ORDER BY id DESC
	`, userID)
	if err != nil {
		return nil, fmt.Errorf("error querying orders: %w", err)
	}
	defer rows.Close()

	var orders []Order
	for rows.Next() {
		var order Order
		if err := rows.Scan(&order.ID, &order.Total); err != nil {
			log.Printf("Error scanning order: %v", err)
			continue
		}
		order.UserID = userID
		items, err := getOrderItems(order.ID)
		if err != nil {
			log.Printf("Error getting order items: %v", err)
			continue
		}
		order.Items = items
		orders = append(orders, order)
	}
	return orders, nil
}

func addToCart(userID int64, itemID int) error {
	// Fetch current cart
	var itemsStr string
	err := db.QueryRow("SELECT items FROM carts WHERE user_id = $1", userID).Scan(&itemsStr)
	if err != nil && err != sql.ErrNoRows {
		return fmt.Errorf("error fetching cart: %w", err)
	}

	var itemsMap map[int]int
	if itemsStr != "" {
		itemsMap = make(map[int]int)
		pairs := strings.Split(itemsStr, ",")
		for _, pair := range pairs {
			parts := strings.Split(pair, ":")
			if len(parts) != 2 {
				continue
			}
			id, err := strconv.Atoi(parts[0])
			if err != nil {
				continue
			}
			qty, err := strconv.Atoi(parts[1])
			if err != nil {
				continue
			}
			itemsMap[id] = qty
		}
	} else {
		itemsMap = make(map[int]int)
	}

	// Increment quantity
	itemsMap[itemID]++

	// Reconstruct items string
	var newItems []string
	for id, qty := range itemsMap {
		newItems = append(newItems, fmt.Sprintf("%d:%d", id, qty))
	}
	newItemsStr := strings.Join(newItems, ",")

	// Upsert cart
	_, err = db.Exec(`
		INSERT INTO carts (user_id, items)
		VALUES ($1, $2)
		ON CONFLICT (user_id) DO UPDATE SET items = EXCLUDED.items
	`, userID, newItemsStr)
	if err != nil {
		return fmt.Errorf("error updating cart: %w", err)
	}

	return nil
}

func clearCart(userID int64) error {
	_, err := db.Exec("DELETE FROM carts WHERE user_id = $1", userID)
	return err
}

func placeOrder(userID int64) error {
	// Fetch cart
	var itemsStr string
	err := db.QueryRow("SELECT items FROM carts WHERE user_id = $1", userID).Scan(&itemsStr)
	if err != nil {
		if err == sql.ErrNoRows {
			return fmt.Errorf("cart is empty")
		}
		return fmt.Errorf("error fetching cart: %w", err)
	}

	// Parse items
	itemsMap := make(map[int]int)
	pairs := strings.Split(itemsStr, ",")
	for _, pair := range pairs {
		parts := strings.Split(pair, ":")
		if len(parts) != 2 {
			continue
		}
		id, err := strconv.Atoi(parts[0])
		if err != nil {
			continue
		}
		qty, err := strconv.Atoi(parts[1])
		if err != nil {
			continue
		}
		itemsMap[id] = qty
	}

	if len(itemsMap) == 0 {
		return fmt.Errorf("cart is empty")
	}

	// Calculate total
	var total float64
	for id, qty := range itemsMap {
		var price float64
		err := db.QueryRow("SELECT price FROM menu_items WHERE id = $1", id).Scan(&price)
		if err != nil {
			return fmt.Errorf("error fetching price for item %d: %w", id, err)
		}
		total += price * float64(qty)
	}

	// Insert order
	var orderID int
	err = db.QueryRow(`
		INSERT INTO orders (user_id, total)
		VALUES ($1, $2)
		RETURNING id
	`, userID, total).Scan(&orderID)
	if err != nil {
		return fmt.Errorf("error inserting order: %w", err)
	}

	// Insert order items
	for id, qty := range itemsMap {
		_, err := db.Exec(`
			INSERT INTO order_items (order_id, item_id, quantity)
			VALUES ($1, $2, $3)
		`, orderID, id, qty)
		if err != nil {
			return fmt.Errorf("error inserting order item: %w", err)
		}
	}

	// Clear cart
	err = clearCart(userID)
	if err != nil {
		return fmt.Errorf("error clearing cart: %w", err)
	}

	return nil
}

func storeUserInfo(userID int64, phone string, location string) error {
	_, err := db.Exec(`
		INSERT INTO users (user_id, phone, location)
		VALUES ($1, $2, $3)
		ON CONFLICT (user_id) DO UPDATE SET phone = EXCLUDED.phone, location = EXCLUDED.location
	`, userID, phone, location)
	return err
}

// Bot functions
func initBot() {
	var err error
	// Fetch bot token from environment variables
	botToken := "7961675553:AAGPX7myxcS3EBj-Rz7RlSXdrXwXunFbW0Q"
	if botToken == "" {
		log.Fatal("TELEGRAM_BOT_API_TOKEN is not set")
	}

	bot, err = tgbotapi.NewBotAPI(botToken)
	if err != nil {
		log.Fatalf("Failed to create bot: %v", err)
	}

	bot.Debug = true
	log.Printf("Authorized on account %s", bot.Self.UserName)
}

func handleUpdates() {
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := bot.GetUpdatesChan(u)

	for update := range updates {
		if update.Message != nil {
			handleMessage(update.Message)
		} else if update.CallbackQuery != nil {
			handleCallbackQuery(update.CallbackQuery)
		}
	}
}

func handleMessage(message *tgbotapi.Message) {
	if message.Text == "/start" {
		sendWelcomeMessage(message.Chat.ID)
	} else if message.Contact != nil {
		// Handle contact information
		userID := message.Contact.UserID
		if userID != message.From.ID {
			// Prevent spoofing of contact
			msg := tgbotapi.NewMessage(message.Chat.ID, "Invalid contact information.")
			bot.Send(msg)
			return
		}
		phone := message.Contact.PhoneNumber
		err := storeUserInfo(userID, phone, "")
		if err != nil {
			log.Printf("Error storing user contact: %v", err)
			msg := tgbotapi.NewMessage(message.Chat.ID, "Failed to store your contact information.")
			bot.Send(msg)
			return
		}
		msg := tgbotapi.NewMessage(message.Chat.ID, "Thank you for sharing your contact!")
		bot.Send(msg)
	} else if message.Location != nil {
		// Handle location information
		userID := message.From.ID
		latitude := message.Location.Latitude
		longitude := message.Location.Longitude
		location := fmt.Sprintf("%f,%f", latitude, longitude)
		err := storeUserInfo(userID, "", location)
		if err != nil {
			log.Printf("Error storing user location: %v", err)
			msg := tgbotapi.NewMessage(message.Chat.ID, "Failed to store your location.")
			bot.Send(msg)
			return
		}
		msg := tgbotapi.NewMessage(message.Chat.ID, "Thank you for sharing your location!")
		bot.Send(msg)
	} else if strings.HasPrefix(message.Text, "/order") {
		// Handle placing an order
		err := placeOrder(message.From.ID)
		if err != nil {
			msg := tgbotapi.NewMessage(message.Chat.ID, fmt.Sprintf("Failed to place order: %v", err))
			bot.Send(msg)
			return
		}
		msg := tgbotapi.NewMessage(message.Chat.ID, "Your order has been placed successfully!")
		bot.Send(msg)
	} else {
		msg := tgbotapi.NewMessage(message.Chat.ID, "I didn't understand that command. Please use /start to see options.")
		bot.Send(msg)
	}
}

func handleCallbackQuery(query *tgbotapi.CallbackQuery) {
	data := query.Data

	switch {
	case data == "menu":
		sendMenu(query.Message.Chat.ID)
	case data == "main_menu":
		sendMainMenu(query.Message.Chat.ID)
	case data == "my_orders":
		sendOrders(query.Message.Chat.ID, query.From.ID)
	case strings.HasPrefix(data, "add_to_cart_"):
		itemIDStr := strings.TrimPrefix(data, "add_to_cart_")
		itemID, err := strconv.Atoi(itemIDStr)
		if err != nil {
			log.Printf("Invalid item ID: %v", err)
			answerCallback(query.ID, "Invalid item.")
			return
		}

		err = addToCart(query.From.ID, itemID)
		if err != nil {
			log.Printf("Error adding to cart: %v", err)
			answerCallback(query.ID, "Failed to add to cart.")
			return
		}

		answerCallback(query.ID, "Added to cart!")
	case data == "view_cart":
		sendCart(query.Message.Chat.ID, query.From.ID)
	case data == "checkout":
		err := placeOrder(query.From.ID)
		if err != nil {
			answerCallback(query.ID, fmt.Sprintf("Failed to place order: %v", err))
			return
		}
		answerCallback(query.ID, "Your order has been placed!")
		sendMainMenu(query.Message.Chat.ID)
	case data == "clear_cart":
		err := clearCart(query.From.ID)
		if err != nil {
			log.Printf("Error clearing cart: %v", err)
			answerCallback(query.ID, "Failed to clear cart.")
			return
		}
		answerCallback(query.ID, "Your cart has been cleared.")
		sendMainMenu(query.Message.Chat.ID)
	case data == "share_contact":
		// This case should be handled via message.Contact
		answerCallback(query.ID, "Please use the button to share your contact.")
	case data == "share_location":
		// This case should be handled via message.Location
		answerCallback(query.ID, "Please use the button to share your location.")
	default:
		answerCallback(query.ID, "Unknown action.")
	}
}

func answerCallback(callbackID, text string) {
	callback := tgbotapi.CallbackConfig{
		CallbackQueryID: callbackID,
		Text:            text,
		ShowAlert:       false,
	}

	if _, err := bot.Request(callback); err != nil {
		log.Printf("Failed to answer callback query: %v", err)
	}
}

func sendWelcomeMessage(chatID int64) {
	// Create a reply keyboard with buttons to share contact and location
	contactButton := tgbotapi.NewKeyboardButton("Share Phone Number")
	contactButton.RequestContact = true

	locationButton := tgbotapi.NewKeyboardButton("Share Location")
	locationButton.RequestLocation = true

	keyboard := tgbotapi.NewReplyKeyboard(
		tgbotapi.NewKeyboardButtonRow(contactButton),
		tgbotapi.NewKeyboardButtonRow(locationButton),
	)

	msg := tgbotapi.NewMessage(chatID, "Welcome to our restaurant bot! Please share your contact and location to proceed.")
	msg.ReplyMarkup = keyboard

	if _, err := bot.Send(msg); err != nil {
		log.Printf("Failed to send welcome message: %v", err)
	}
}

func sendMainMenu(chatID int64) {
	keyboard := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("📖 Menu", "menu"),
			tgbotapi.NewInlineKeyboardButtonData("🛒 My Cart", "view_cart"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("📦 My Orders", "my_orders"),
		),
	)

	msg := tgbotapi.NewMessage(chatID, "Main Menu: Please choose an option:")
	msg.ReplyMarkup = keyboard
	if _, err := bot.Send(msg); err != nil {
		log.Printf("Failed to send main menu: %v", err)
	}
}

func sendMenu(chatID int64) {
	items, err := getMenuItems()
	if err != nil {
		log.Printf("Error getting menu items: %v", err)
		msg := tgbotapi.NewMessage(chatID, "Failed to retrieve menu.")
		bot.Send(msg)
		return
	}

	var keyboard [][]tgbotapi.InlineKeyboardButton
	for _, item := range items {
		button := tgbotapi.NewInlineKeyboardButtonData(
			fmt.Sprintf("%s - %.0f so'm", item.Name, item.Price),
			fmt.Sprintf("add_to_cart_%d", item.ID),
		)
		row := tgbotapi.NewInlineKeyboardRow(button)
		keyboard = append(keyboard, row)
	}

	// Add a back button
	backButton := tgbotapi.NewInlineKeyboardButtonData("🔙 Back to Main Menu", "main_menu")
	keyboard = append(keyboard, tgbotapi.NewInlineKeyboardRow(backButton))

	msg := tgbotapi.NewMessage(chatID, "Our Menu:")
	msg.ReplyMarkup = tgbotapi.NewInlineKeyboardMarkup(keyboard...)
	if _, err := bot.Send(msg); err != nil {
		log.Printf("Failed to send menu: %v", err)
	}
}

func sendOrders(chatID int64, userID int64) {
	orders, err := getUserOrders(userID)
	if err != nil {
		log.Printf("Error getting user orders: %v", err)
		msg := tgbotapi.NewMessage(chatID, "Failed to retrieve your orders.")
		bot.Send(msg)
		return
	}

	if len(orders) == 0 {
		msg := tgbotapi.NewMessage(chatID, "You have no orders yet.")
		bot.Send(msg)
		return
	}

	var response strings.Builder
	for _, order := range orders {
		response.WriteString(fmt.Sprintf("🆔 Order #%d\n💰 Total: %.0f so'm\n\n", order.ID, order.Total))
	}

	msg := tgbotapi.NewMessage(chatID, response.String())
	bot.Send(msg)
}

func sendCart(chatID int64, userID int64) {
	// Fetch cart
	var itemsStr string
	err := db.QueryRow("SELECT items FROM carts WHERE user_id = $1", userID).Scan(&itemsStr)
	if err != nil {
		if err == sql.ErrNoRows {
			msg := tgbotapi.NewMessage(chatID, "Your cart is empty.")
			bot.Send(msg)
			return
		}
		log.Printf("Error fetching cart: %v", err)
		msg := tgbotapi.NewMessage(chatID, "Failed to retrieve your cart.")
		bot.Send(msg)
		return
	}

	// Parse items
	itemsMap := make(map[int]int)
	pairs := strings.Split(itemsStr, ",")
	for _, pair := range pairs {
		parts := strings.Split(pair, ":")
		if len(parts) != 2 {
			continue
		}
		id, err := strconv.Atoi(parts[0])
		if err != nil {
			continue
		}
		qty, err := strconv.Atoi(parts[1])
		if err != nil {
			continue
		}
		itemsMap[id] = qty
	}

	if len(itemsMap) == 0 {
		msg := tgbotapi.NewMessage(chatID, "Your cart is empty.")
		bot.Send(msg)
		return
	}

	// Fetch item details
	var response strings.Builder
	var total float64
	for id, qty := range itemsMap {
		var name string
		var price float64
		err := db.QueryRow("SELECT name, price FROM menu_items WHERE id = $1", id).Scan(&name, &price)
		if err != nil {
			log.Printf("Error fetching item %d: %v", id, err)
			continue
		}
		response.WriteString(fmt.Sprintf("• %s x%d - %.0f so'm\n", name, qty, price*float64(qty)))
		total += price * float64(qty)
	}

	response.WriteString(fmt.Sprintf("\n💰 Total: %.0f so'm", total))

	// Add buttons for checkout, clear cart, and back
	keyboard := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("✅ Checkout", "checkout"),
			tgbotapi.NewInlineKeyboardButtonData("🗑️ Clear Cart", "clear_cart"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("🔙 Back to Main Menu", "main_menu"),
		),
	)

	msg := tgbotapi.NewMessage(chatID, response.String())
	msg.ReplyMarkup = keyboard
	if _, err := bot.Send(msg); err != nil {
		log.Printf("Failed to send cart: %v", err)
	}
}

// Main function
func main() {
	// Initialize logging
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	// Initialize database
	initDB()
	defer db.Close()

	// Initialize bot
	initBot()

	// Start handling updates
	handleUpdates()
}
