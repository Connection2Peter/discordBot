package main


import (
	"os"
	"fmt"
	"time"
	"strings"
	"io/ioutil"
	"math/rand"
	"path/filepath"
	"encoding/json"

	"github.com/bwmarrin/discordgo"
)



// ### User setting
const TOKEN = ""
const USAGE = `
***### Usage Of This Discord Bot*** __**Version 1.2**__

***>usage***
	*Usage of this bot.*
***>change***
	*Change log of this bot.*
***>image***
	*Get one random Namin photo. (PNG/ JPG/ GIF)*
***>food***
	*Get one random recommendation of cuisine.*
***>foodImage***
	*Get one random image of cuisine.*
***>food2Hellaya***
	*Get one random image of cuisine. (0.01% King Crab)*
***>gif***
	*Get one random Namin photo. (only GIF)*

***Most images with Namin powered by "ArChi#6749".***
***If there's an image is not good-looking,***
***or this bot with any problem.***
***Please mention "connection#8506" to fix them up.***
***thanks for your use !!!***
`
const CHANGE = `
***### Change Log*** __**Version 1.2**__
***# 1. Hive off gif from image.***
`


// ### Global vars & type
type Food struct {
	ID          int64    `json:"id"`         
	Cuisine     string   `json:"cuisine"`    
	Ingredients []string `json:"ingredients"`
}

var KY int = 0
var LD int = 0
var BOTID string
var ImgNamin []string
var ImgNmGif []string
var ImgFoods []string
var FOODS []Food



// ### Event
func messageHandler(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.ID == BOTID {
		return
	}

	inputText := strings.ReplaceAll(strings.ToLower(m.Content), " ", "")

	if strings.Contains(inputText, ">usage") {
		s.ChannelMessageSend(m.ChannelID, USAGE)
	}

	if strings.Contains(inputText, ">change") {
		s.ChannelMessageSend(m.ChannelID, CHANGE)
	}

	if m.Content == ">image" {
		fcName := ImgNamin[rand.Int() % len(ImgNamin)]

		f, err := os.Open(fcName)
		if err != nil { fmt.Println("os.Open", err) }
		defer f.Close()

		msg := "<@" + m.Author.ID + ">\n"

		_, err = s.ChannelFileSendWithMessage(m.ChannelID, msg, filepath.Base(fcName), f)
		if err != nil {
			fmt.Println(err)
		}
	}

	if m.Content == ">gif" {
		fcName := ImgNmGif[rand.Int() % len(ImgNmGif)]

		f, err := os.Open(fcName)
		if err != nil { fmt.Println("os.Open", err) }
		defer f.Close()

		msg := "<@" + m.Author.ID + ">\n"

		_, err = s.ChannelFileSendWithMessage(m.ChannelID, msg, filepath.Base(fcName), f)
		if err != nil {
			fmt.Println(err)
		}
	}

	if m.Content == ">food" {
		food := FOODS[rand.Int() % len(FOODS)]
		text := "<@" + m.Author.ID + ">\n**Food Type** : " + food.Cuisine + "\n**Ingredients** : \n"

		s.ChannelMessageSend(m.ChannelID, text + strings.Join(food.Ingredients, " "))
	}

	if m.Content ==  ">foodImage" {
		fcName := ImgFoods[rand.Int() % len(ImgFoods)]

		f, err := os.Open(fcName)
		if err != nil { fmt.Println("os.Open", err) }
		defer f.Close()

		msg := "<@" + m.Author.ID + ">\n"

		_, err = s.ChannelFileSendWithMessage(m.ChannelID, msg, filepath.Base(fcName), f)
		if err != nil {
			fmt.Println(err)
		}
	}

	if m.Content == ">food2Hellaya" {
		fcName := ImgFoods[rand.Int() % len(ImgFoods)]

		f, err := os.Open(fcName)
		if err != nil { fmt.Println("os.Open", err) }
		defer f.Close()

		msg := "***Good morning Alan :***\n"

		_, err = s.ChannelFileSendWithMessage(m.ChannelID, msg, filepath.Base(fcName), f)
		if err != nil {
			fmt.Println(err)
		}
	}

	if strings.Contains(inputText, "幹你媽ky") {
		KY++
		msg := fmt.Sprintf("KY已經被幹了 %d 次", KY)

		s.ChannelMessageSend(m.ChannelID, msg)
	}

	if strings.Contains(inputText, "幹你媽衣服") {
		LD++
		msg := fmt.Sprintf("衣服已經被幹了 %d 次", LD)

		s.ChannelMessageSend(m.ChannelID, msg)
	}
}



// ### Main
func main() {
	RequireArgs := []string{
		"pathFoodJson",
		"pathImgNamin",
		"pathImgFoods",
	}

	if len(os.Args) < len(RequireArgs)+1 {
		usage := "go run " + os.Args[0]

		for _, arg := range RequireArgs {
			usage += " " + arg
		}

		fmt.Printf("Usage :\n %s\n", usage)

		return
	} else {
		fmt.Println("### Argument List :")
		fmt.Println("# This Program\t\t:", os.Args[0])

		for idx, arg := range RequireArgs {
			fmt.Printf("# %s\t: %s\n", arg, os.Args[idx+1])
		}
	}

	goBot, err := discordgo.New("Bot " + TOKEN)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	u, err := goBot.User("@me")
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	ImgNamin, err = filepath.Glob(filepath.Join(os.Args[2], "*"))
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	ImgNmGif, err = filepath.Glob(filepath.Join(os.Args[2], "*.gif"))
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	ImgFoods, err = filepath.Glob(filepath.Join(os.Args[3], "*"))
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	jsonFile, err := os.Open(os.Args[1])
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	defer jsonFile.Close()

	jsonBuf, _ := ioutil.ReadAll(jsonFile)

	err = json.Unmarshal(jsonBuf, &FOODS)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	rand.Seed(time.Now().UnixNano())

	BOTID = u.ID

	goBot.AddHandler(messageHandler)

	err = goBot.Open()
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	fmt.Println("### Bot Online.")

	<-make(chan struct{})
}
