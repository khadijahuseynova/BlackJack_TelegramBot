package main

import (
	"log"
	"math/rand"
	"strconv"
	"strings"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

type Card struct {
	Value int
	Suit  string
}

type Hand struct {
	Cards []Card
}

type Game struct {
	PlayerHand Hand
	DealerHand Hand
}

var game *Game

var suits = []string{"♥️", "♦️", "♣️", "♠️"}

var values = []int{2, 3, 4, 5, 6, 7, 8, 9, 10, 10, 10, 10, 11}

func main() {
	rand.Seed(time.Now().UnixNano())

	bot, err := tgbotapi.NewBotAPI("6210283770:AAEEj7MS4zr6KW8vBMTRsGqOckHolK4AdMk") 
	if err != nil {
		log.Fatal(err)
	}

	bot.Debug = true

	log.Printf("Authorized on account %s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates, err := bot.GetUpdatesChan(u)
	if err != nil {
		log.Fatal(err)
	}

	for update := range updates {
		
		if update.Message == nil {
			continue
		}

		if update.Message.Text == "/start" {
			game = NewGame()
			
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Blackjack oyununa hoş geldiniz!\n\nOyunu başlatmak için /play komutunu kullanın.")
			
			bot.Send(msg)
		} else if update.Message.Text == "/play" {
			if game != nil {
			
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Oyun zaten başladı!")
				
				bot.Send(msg)
			} else {
				game = NewGame()
			
			
				game.DealInitialCards()
				
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Kartlar dağıtıldı!\n\nElinizdeki kartlar: "+game.PlayerHand.String()+"\n\nDeğer: "+strconv.Itoa(game.PlayerHand.Value()))
				
				bot.Send(msg)
			}

		} else if update.Message.Text == "hit" {
			if game != nil && game.IsInProgress() {
				card := game.DealCard()
				game.PlayerHand.Cards = append(game.PlayerHand.Cards, card)



				msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Yeni kart: "+card.String()+"\n\nElinizdeki kartlar: "+game.PlayerHand.String()+"\n\nDeğer: "+strconv.Itoa(game.PlayerHand.Value()))
				bot.Send(msg)

				if game.PlayerHand.Value() > 21 {
					msg = tgbotapi.NewMessage(update.Message.Chat.ID, " 21'i aştı! Oyunu kaybettiniz:(")
					bot.Send(msg)
					game = nil
				} else if game.PlayerHand.Value() == 21 {
					msg = tgbotapi.NewMessage(update.Message.Chat.ID, "BlackJack!!!")
					bot.Send(msg)
					game = nil
				}
			} else {
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Oyun henüz başlamadı! /play komutuyla başlatın.")
				bot.Send(msg)
			}
		} else if update.Message.Text == "stand" {
			if game != nil && game.IsInProgress() {
				for game.DealerHand.Value() < 17 {
					card := game.DealCard()
					game.DealerHand.Cards = append(game.DealerHand.Cards, card)
				}

				msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Dağıtıcının kartları: "+game.DealerHand.String()+"\n\nDeğer: "+strconv.Itoa(game.DealerHand.Value()))
				
				
				bot.Send(msg)

				
				
				if game.DealerHand.Value() > 21 || game.PlayerHand.Value() > game.DealerHand.Value() {
					msg = tgbotapi.NewMessage(update.Message.Chat.ID, "Tebrikler! Kazandınız!!!")
					bot.Send(msg)
				} else if game.DealerHand.Value() > game.PlayerHand.Value() {
					msg = tgbotapi.NewMessage(update.Message.Chat.ID, "Üzgünüm, oyunu kaybettiniz:(")
					bot.Send(msg)
				} else {
					msg = tgbotapi.NewMessage(update.Message.Chat.ID, "Oyun berabere sonuçland:|")
					bot.Send(msg)
				}

				game = nil
			} else {
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Oyun henüz başlamadı! /play komutuyla başlatın.")
				bot.Send(msg)
			}
		}
	}
}

func NewGame() *Game {
	return &Game{
		PlayerHand: Hand{},
		
		DealerHand: Hand{},
	}
}

func (g *Game) DealInitialCards() {



	g.PlayerHand.Cards = append(g.PlayerHand.Cards, g.DealCard(), g.DealCard())
	
	g.DealerHand.Cards = append(g.DealerHand.Cards, g.DealCard())
}

func (g *Game) DealCard() Card {
	card := Card{
		Value: values[rand.Intn(len(values))],
		
		
		Suit:  suits[rand.Intn(len(suits))],
	}
	return card
}

func (h Hand) Value() int {
	value := 0
	aceCount := 0

	for _, card := range h.Cards {
		if card.Value == 11 {
			aceCount++
		}
		value += card.Value
	}





	for i := 0; i < aceCount; i++ {
		if value > 21 {
			value -= 10
		}
	}

	return value
}

func (h Hand) String() string {
	var cardStrings []string
	for _, card := range h.Cards {
		cardStrings = append(cardStrings, card.String())
	}
	return strings.Join(cardStrings, ", ")
}


func (c Card) String() string {
	
	return strconv.Itoa(c.Value) + " " + c.Suit


}

func (g *Game) IsInProgress() bool {
	
	return len(g.PlayerHand.Cards) > 0 && len(g.DealerHand.Cards) > 0
}
