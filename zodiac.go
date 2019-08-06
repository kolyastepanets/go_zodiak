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
  "net/http"
  "strings"
  "github.com/go-redis/redis"
  "time"
)

var (
  clientRedis *redis.Client
  err error
)

const objectPrefix string = "user_id_"
const maxQuantityToGenerateRudePart1 int = 5
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

// init is invoked before main()
func init() {
  // loads values from .env into the system
  if err := godotenv.Load(); err != nil {
    log.Print("No .env file found")
  }

  redisAddr, _ := os.LookupEnv("REDIS_ADDR")
  redisPwd, _ := os.LookupEnv("REDIS_PWD")

  clientRedis = redis.NewClient(&redis.Options{
    Addr:     redisAddr,
    Password: redisPwd, // no password set
    DB:       0,  // use default DB
  })

  pong, err := clientRedis.Ping().Result()
  fmt.Println(pong, err)
}

func MainHandler(resp http.ResponseWriter, _ *http.Request) {
  resp.Write([]byte("Hi there!"))
}

func CanGenerateHoroscope(update tgbotapi.Update) bool {
  UserObj := GetUser(update)

  var isDatesEqual bool
  isDatesEqual = DateEqual(UserObj.DateRequestHoroscope, time.Now())
  return !isDatesEqual
}

func GenerateAbuseMessage(update tgbotapi.Update) string {
  UserObj := GetUser(update)

  rude_reply, err := ioutil.ReadFile("rude_reply.json")
  if err != nil {
    fmt.Print(err)
  }

  var RudeReplyObj RudeReply

  err = json.Unmarshal([]byte(rude_reply), &RudeReplyObj)
  if err != nil {
    fmt.Println("error:", err)
  }

  var sentence string

  if UserObj.TimesHoroscopeWasRequested >= maxQuantityToGenerateRudePart1 {
    var index = UserObj.TimesHoroscopeWasRequested - maxQuantityToGenerateRudePart1
    if index >= len(RudeReplyObj.Part2) {
      index = RandomNumber(len(RudeReplyObj.Part2) - 1)
    }
    sentence = RudeReplyObj.Part2[index]
  } else {
    sentence = RudeReplyObj.Part1[UserObj.TimesHoroscopeWasRequested]
  }

  UserObj.TimesHoroscopeWasRequested = UserObj.TimesHoroscopeWasRequested + 1
  SaveUser(UserObj)

  return sentence
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

  var sentence = currentZodiakSign.Index(RandomNumber(currentZodiakSign.Len())).Interface().(string)
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

func RandomNumber(max int) int {
  return rand.Intn((max - 1))
}

func DateEqual(date1, date2 time.Time) bool {
  y1, m1, d1 := date1.Date()
  y2, m2, d2 := date2.Date()

  return y1 == y2 && m1 == m2 && d1 == d2
}

func GenerateHoroscope(update tgbotapi.Update) string {
  horoscope, err := ioutil.ReadFile("horoscope_generator.json")
  if err != nil {
    fmt.Print(err)
  }

  var results [][]interface{}

  err = json.Unmarshal([]byte(horoscope), &results)
  if err != nil {
    fmt.Println("error:", err)
  }
  var sentence []string

  for _, result := range results {
    if str, ok := result[RandomNumber(len(result))].(string); ok {
      sentence = append(sentence, str)
    }
  }

  UserObj := GetUser(update)
  UserObj.TimesHoroscopeWasRequested = 0
  UserObj.DateRequestHoroscope = time.Now()
  SaveUser(UserObj)

  return strings.Join(sentence, "")
}

func CallbackHandler(callback tgbotapi.CallbackQuery, bot *tgbotapi.BotAPI) {
  var horoscope = "\nГороскоп на сегодня, кликни /horoscope"
  var sentence = FindRussianNameForZodiak(callback.Data) + ": " + FindSentenceForZodiac(callback.Data) + horoscope

  msg := tgbotapi.NewMessage(callback.Message.Chat.ID, callback.Message.Text)
  msg.Text = sentence
  msg.ReplyMarkup = zodiacKeyboard
  bot.Send(msg)
}

func main() {
  http.HandleFunc("/", MainHandler)
  go http.ListenAndServe(":"+os.Getenv("PORT"), nil)

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

  updates := bot.ListenForWebhook("/" + telegramBotToken)

  for update := range updates {
    if update.CallbackQuery != nil {
      CallbackHandler(*update.CallbackQuery, bot)
    } else if update.Message.IsCommand() {
      msg := tgbotapi.NewMessage(update.Message.Chat.ID, update.Message.Text)
      switch update.Message.Command() {
      case "zodiac":
        msg.Text = "Выбери знак зодиака"
        msg.ReplyMarkup = zodiacKeyboard
      case "help":
        msg.Text = "кликни /start"
      case "horoscope":
        if CanGenerateHoroscope(update) {
          msg.Text = GenerateHoroscope(update)
        } else {
          msg.Text = GenerateAbuseMessage(update)
        }
      case "start":
        msg.Text = "Очень точное описание знаков зодиака, (осторожно мат), кликни /zodiac"
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
