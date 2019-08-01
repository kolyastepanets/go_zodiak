package main

import (
  "log"
  "fmt"
  "github.com/go-telegram-bot-api/telegram-bot-api"
  "io/ioutil"
  "encoding/json"
  "reflect"
  "math/rand"
  "github.com/joho/godotenv"
  "os"
)

var zodiacKeyboard = tgbotapi.NewInlineKeyboardMarkup(
  tgbotapi.NewInlineKeyboardRow(
    tgbotapi.NewInlineKeyboardButtonData("Овен", "Aries"),
    tgbotapi.NewInlineKeyboardButtonData("Телец", "Taurus"),
    tgbotapi.NewInlineKeyboardButtonData("Близнецы", "Gemini"),
  ),
  tgbotapi.NewInlineKeyboardRow(
    tgbotapi.NewInlineKeyboardButtonData("Рак", "Cancer"),
    tgbotapi.NewInlineKeyboardButtonData("Лев", "Leo"),
    tgbotapi.NewInlineKeyboardButtonData("Дева", "Virgo"),
  ),
  tgbotapi.NewInlineKeyboardRow(
    tgbotapi.NewInlineKeyboardButtonData("Весы", "Libra"),
    tgbotapi.NewInlineKeyboardButtonData("Скорпион", "Scorpio"),
    tgbotapi.NewInlineKeyboardButtonData("Стрелец", "Saggitarius"),
  ),
  tgbotapi.NewInlineKeyboardRow(
    tgbotapi.NewInlineKeyboardButtonData("Козерог", "Capricorn"),
    tgbotapi.NewInlineKeyboardButtonData("Водолей", "Aquarius"),
    tgbotapi.NewInlineKeyboardButtonData("Рыбы", "Pisces"),
  ),
)

type ZodiacSigns struct {
  Aquarius []string `json:"aquarius"`
  Pisces []string `json:"pisces"`
  Aries []string `json:"aries"`
  Taurus []string `json:"taurus"`
  Gemini []string `json:"gemini"`
  Cancer []string `json:"cancer"`
  Leo []string `json:"leo"`
  Virgo []string `json:"virgo"`
  Libra []string `json:"libra"`
  Scorpio []string `json:"scorpio"`
  Saggitarius []string `json:"saggitarius"`
  Capricorn []string `json:"capricorn"`
}

type RussianZodiacs struct {
  Aquarius string `json:"aquarius"`
  Pisces string `json:"pisces"`
  Aries string `json:"aries"`
  Taurus string `json:"taurus"`
  Gemini string `json:"gemini"`
  Cancer string `json:"cancer"`
  Leo string `json:"leo"`
  Virgo string `json:"virgo"`
  Libra string `json:"libra"`
  Scorpio string `json:"scorpio"`
  Saggitarius string `json:"saggitarius"`
  Capricorn string `json:"capricorn"`
}

// init is invoked before main()
func init() {
  // loads values from .env into the system
  if err := godotenv.Load(); err != nil {
    log.Print("No .env file found")
  }
}

func FindSentenceForZodiac(callbackData string) string {
  zodiacSigns, err := ioutil.ReadFile("zodiac_signs.json")
  if err != nil {
    fmt.Print(err)
  }

  var ZodiacSignsObj ZodiacSigns

  err = json.Unmarshal([]byte(zodiacSigns), &ZodiacSignsObj)
  if err != nil {
    fmt.Println("error:", err)
  }

  reflectZodiacSigns := reflect.ValueOf(ZodiacSignsObj)
  currentZodiakSign := reflect.Indirect(reflectZodiacSigns).FieldByName(callbackData)

  var sentence = currentZodiakSign.Index(rand.Intn((currentZodiakSign.Len() - 1) - 0)).Interface().(string)
  return sentence
}

func FindRussianNameForZodiak(callbackData string) string {
  russianZodiacs, err := ioutil.ReadFile("russian_zodiac.json")
  if err != nil {
    fmt.Print(err)
  }

  var RussianZodiacsObj RussianZodiacs

  err = json.Unmarshal([]byte(russianZodiacs), &RussianZodiacsObj)
  if err != nil {
    fmt.Println("error:", err)
  }

  reflectRussianZodiacSigns := reflect.ValueOf(RussianZodiacsObj)
  currentRussianZodiacSign := reflect.Indirect(reflectRussianZodiacSigns).FieldByName(callbackData)

  return currentRussianZodiacSign.Interface().(string)
}

func CallbackHandler(callback tgbotapi.CallbackQuery, bot *tgbotapi.BotAPI) {
  var sentence = FindRussianNameForZodiak(callback.Data) + ": " + FindSentenceForZodiac(callback.Data)

  msg := tgbotapi.NewMessage(callback.Message.Chat.ID, callback.Message.Text)
  msg.Text = sentence
  msg.ReplyMarkup = zodiacKeyboard
  bot.Send(msg)
}

func main() {
  telegramBotToken, exists := os.LookupEnv("TELEGRAM_BOT_TOKEN")

  if exists {
    fmt.Println("token exists")
  }

  bot, err := tgbotapi.NewBotAPI(telegramBotToken)
  if err != nil {
    log.Panic(err)
  }

  // bot.Debug = true

  // log.Printf("Authorized on account %s", bot.Self.UserName)

  u := tgbotapi.NewUpdate(0)
  u.Timeout = 60

  updates, err := bot.GetUpdatesChan(u)

  for update := range updates {
    if update.CallbackQuery != nil {
      CallbackHandler(*update.CallbackQuery, bot)
    } else if update.Message.IsCommand() {
      msg := tgbotapi.NewMessage(update.Message.Chat.ID, update.Message.Text)
      switch update.Message.Command() {
      case "choose":
        msg.Text = "Выбери знак зодиака"
        msg.ReplyMarkup = zodiacKeyboard
      case "help":
        msg.Text = "кликни /start"
      case "start":
        msg.Text = "Очень точное описание знаков зодиака, (осторожно мат), кликни /choose"
      default:
        msg.Text = "Нет такой команды"
      }
      bot.Send(msg)
    } else if update.Message != nil {
      msg := tgbotapi.NewMessage(update.Message.Chat.ID, update.Message.Text)
      msg.Text = "Нет такой команды"
      bot.Send(msg)
    }
  }
}
