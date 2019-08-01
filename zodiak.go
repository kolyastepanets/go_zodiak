package main

import (
  "log"
  "fmt"
  "github.com/go-telegram-bot-api/telegram-bot-api"
  "io/ioutil"
  "encoding/json"
  "reflect"
  "math/rand"
)

var zodiakKeyboard = tgbotapi.NewInlineKeyboardMarkup(
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

type ZodiakSigns struct {
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

type RussianZodiaks struct {
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

func FindSentenceForZodiak(callbackData string) string {
  zodiakSigns, err := ioutil.ReadFile("zodiak_signs.json")
  if err != nil {
    fmt.Print(err)
  }

  var ZodiakSignsObj ZodiakSigns

  err = json.Unmarshal([]byte(zodiakSigns), &ZodiakSignsObj)
  if err != nil {
    fmt.Println("error:", err)
  }

  reflectZodiakSigns := reflect.ValueOf(ZodiakSignsObj)
  currentZodiakSign := reflect.Indirect(reflectZodiakSigns).FieldByName(callbackData)

  var sentence = currentZodiakSign.Index(rand.Intn((currentZodiakSign.Len() - 1) - 0)).Interface().(string)
  return sentence
}

func FindRussianNameForZodiak(callbackData string) string {
  russianZodiaks, err := ioutil.ReadFile("russian_zodiak.json")
  if err != nil {
    fmt.Print(err)
  }

  var RussianZodiaksObj RussianZodiaks

  err = json.Unmarshal([]byte(russianZodiaks), &RussianZodiaksObj)
  if err != nil {
    fmt.Println("error:", err)
  }

  reflectRussianZodiakSigns := reflect.ValueOf(RussianZodiaksObj)
  currentRussianZodiakSign := reflect.Indirect(reflectRussianZodiakSigns).FieldByName(callbackData)

  return currentRussianZodiakSign.Interface().(string)
}

func CallbackHandler(callback tgbotapi.CallbackQuery, bot *tgbotapi.BotAPI) {
  var sentence = FindRussianNameForZodiak(callback.Data) + ": " + FindSentenceForZodiak(callback.Data)

  msg := tgbotapi.NewMessage(callback.Message.Chat.ID, callback.Message.Text)
  msg.Text = sentence
  msg.ReplyMarkup = zodiakKeyboard
  bot.Send(msg)
}

func main() {
  bot, err := tgbotapi.NewBotAPI("610859316:AAHO5IY_npP8Bszm_1oQW_vPf7myqu30vYw")
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
        msg.ReplyMarkup = zodiakKeyboard
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
