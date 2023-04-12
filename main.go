package main

import (
	"encoding/csv"
	"fmt"
	"os"
	"strconv"
	"time"
)

type command func(app *Application) error

type Application struct {
	FileName string
	Commands map[string]command
}

type Godo struct {
	Id      int
	Date    string
	Message string
}

func main() {
	home, err := os.UserHomeDir()
	if err != nil {
		fmt.Println("Error: " + err.Error())
		os.Exit(1)
	}

	app := NewApplication(home + "/godos.csv")

	app.Cmd("help", HelpCommand)
	app.Cmd("new", NewCommand)
	app.Cmd("all", AllCommand)
	app.Cmd("rm", RemoveCommand)

	if err := app.Run(); err != nil {
		fmt.Println("Error: " + err.Error())
		os.Exit(1)
	}
}

func NewApplication(fileName string) Application {
	return Application{
		FileName: fileName,
		Commands: map[string]command{},
	}
}

func (app *Application) Cmd(name string, cmd command) {
	app.Commands[name] = cmd
}

func (app *Application) Run() error {
	if len(os.Args) > 1 {
		cmd := os.Args[1]

		callback, ok := app.Commands[cmd]
		if ok == false {
			return fmt.Errorf("A command with the name '%s' does not exist", cmd)
		}

		return callback(app)
	}

	return AllCommand(app)
}

func HelpCommand(app *Application) error {
	fmt.Print("Hello and welcome to your CLI todo app 🚀 \n\n")

	fmt.Printf("help\t\tDisplays this help\t\tExample: 'godo help'\n")
	fmt.Printf("new\t\tCreates a new godo\t\tExample: 'godo new \"Clean kitchen\"'\n")
	fmt.Printf("all\t\tLists all current godos\t\tExample: 'godo all'\n")
	fmt.Printf("rm\t\tDeletes a godo\t\t\tExample: 'godo rm 123'\n")

	return nil
}

func NewCommand(app *Application) error {
	godos, err := readCsvFile(app.FileName)
	if err != nil {
		return err
	}

	if len(os.Args) < 3 {
		return fmt.Errorf("Looks like you did not pass a new godo message 🤔")
	}

	id := getNewId(godos)
	godo := Godo{Id: id, Date: time.Now().Format("2006-01-02 15:04:05"), Message: os.Args[2]}
	godos = append(godos, godo)

	if err := writeCsvFile(app.FileName, godos); err != nil {
		return err
	}

	fmt.Printf("Successfully added new godo entry '%s' 💪\n", godo.Message)

	return nil
}

func AllCommand(app *Application) error {
	godos, err := readCsvFile(app.FileName)
	if err != nil {
		return err
	}

	if len(godos) == 0 {
		fmt.Println("No godo entries available 🤷")
	} else {
		for _, godo := range godos {
			fmt.Printf("%d\t%s\t%s\n", godo.Id, godo.Date, godo.Message)
		}
	}

	return nil
}

func RemoveCommand(app *Application) error {
	godos, err := readCsvFile(app.FileName)
	if err != nil {
		return err
	}

	if len(os.Args) < 3 {
		return fmt.Errorf("Looks like you did not pass a godo id 🤔")
	}

	id, err := strconv.Atoi(os.Args[2])
	if err != nil {
		return err
	}

	found := false
	for i, godo := range godos {
		if id == godo.Id {
			found = true
			godos = append(godos[:i], godos[i+1:]...)
			break
		}
	}

	if err := writeCsvFile(app.FileName, godos); err != nil {
		return err
	}

	if found {
		fmt.Printf("Successfully removed godo entry with id '%d' 👌\n", id)
	} else {
		fmt.Printf("No godo entry with id '%d' found 🤷\n", id)
	}

	return nil
}

func readCsvFile(fileName string) ([]Godo, error) {
	f, err := os.OpenFile(fileName, os.O_RDONLY|os.O_CREATE, 0644)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	lines, err := csv.NewReader(f).ReadAll()
	if err != nil {
		return nil, err
	}

	var godos []Godo
	for _, line := range lines {
		id, err := strconv.Atoi(line[0])
		if err != nil {
			return nil, err
		}

		godo := Godo{
			Id:      id,
			Message: line[1],
			Date:    line[2],
		}

		godos = append(godos, godo)
	}

	return godos, nil
}

func writeCsvFile(fileName string, godos []Godo) error {
	f, err := os.OpenFile(fileName, os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		return err
	}
	defer f.Close()

	err = f.Truncate(0)
	if err != nil {
		return err
	}

	w := csv.NewWriter(f)
	for _, godo := range godos {
		if err := w.Write([]string{strconv.Itoa(godo.Id), godo.Message, godo.Date}); err != nil {
			return err
		}
	}
	w.Flush()

	return nil
}

func getNewId(godos []Godo) int {
	id := 1

	for {
		found := false
		for _, godo := range godos {
			if id == godo.Id {
				found = true
				break
			}
		}

		if found == false {
			break
		}

		id++
	}

	return id
}
