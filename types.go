package main

import "time"

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

type User struct {
  Id int64 `json:"id"`
  FirstName string `json:"first_name"`
  LastName string `json:"last_name"`
  DateRequestHoroscope time.Time `json:"date_request_horoscope"`
  TimesHoroscopeWasRequested int `json:"times_horoscope_was_requested"`
}

type RudeReply struct {
  Part1 []string `json:"part_1"`
  Part2 []string `json:"part_2"`
}
