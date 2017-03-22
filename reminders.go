package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"time"
)

type Name struct {
	Date        string
	Description string
}

var (
	newFlag   = flag.Bool("n", false, "create new reminder")
	clearFlag = flag.Bool("c", false, "clear reminders")
	listFlag  = flag.Bool("l", false, "list all reminders")
	helpFlag  = flag.Bool("h", false, "help")
)

func checkFile() {
	if _, err := os.Stat("reminders.json"); os.IsNotExist(err) {
		os.Create("reminders.json")
	}
}

func newReminder() {
	checkFile()

	var text string
	date := time.Now()
	dateString := date.String()

	file, err := os.OpenFile("reminders.json", os.O_RDWR|os.O_APPEND, os.ModePerm)
	if err != nil {
		log.Fatal(err)
	}

	defer file.Close()

	fmt.Print("Description: ")

	scanner := bufio.NewScanner(os.Stdin)
	if ok := scanner.Scan(); ok {
		text = scanner.Text()
	}

	reminder := &Name{dateString[:19], text}
	newReminder, err := json.MarshalIndent(&reminder, "", "    ")
	if err != nil {
		log.Fatal(err)
	}

	file.Write(newReminder)
}

func removeReminders() {
	err := os.Remove("reminders.json")
	if err != nil {
		log.Fatal(err)
	}
}

func listReminders() error {
	file, err := os.Open("reminders.json")
	if err != nil {
		return err
	}

	decoder := json.NewDecoder(file)

	for {
		var content Name
		switch decoder.Decode(&content) {
		case nil:
			fmt.Printf("%s: %s\n", content.Date, content.Description)
		case io.EOF:
			return nil
		default:
			return err
		}
	}
}

func main() {
	flag.Parse()

	if len(os.Args) < 2 {
		fmt.Println("Usage:")
		flag.PrintDefaults()
	}

	if *newFlag {
		newReminder()
	}

	if *clearFlag {
		removeReminders()
	}

	if *helpFlag {
		flag.PrintDefaults()
	}

	if *listFlag {
		err := listReminders()
		if err != nil {
			log.Fatal(err)
		}
	}

}
