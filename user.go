package main

import (
  "fmt"
  "github.com/go-telegram-bot-api/telegram-bot-api"
  "encoding/json"
  "github.com/go-redis/redis"
)

func InitUser(update tgbotapi.Update) {
  usr := User{
    Id:                         update.Message.Chat.ID,
    FirstName:                  update.Message.Chat.FirstName,
    LastName:                   update.Message.Chat.LastName,
    TimesHoroscopeWasRequested: 0,
  }
  SaveUser(usr)
}

func GetUser(update tgbotapi.Update) User {
  user, err := clientRedis.Get(objectPrefix + fmt.Sprint(update.Message.Chat.ID)).Result()
  if err == redis.Nil {
    InitUser(update)
  } else if err != nil {
    panic(err)
  }

  var UserObj User
  err = json.Unmarshal([]byte(user), &UserObj)
  if err != nil {
    fmt.Println("error:", err)
  }

  return UserObj
}

func SaveUser(UserObj User) {
  userJson, err := json.Marshal(UserObj)
  if err != nil {
    fmt.Println("error:", err)
  }

  errRedis := clientRedis.Set(objectPrefix + fmt.Sprint(UserObj.Id), userJson, 0).Err()
  if errRedis != nil {
    panic(errRedis)
  }
}
