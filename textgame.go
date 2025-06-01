package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"

	"gopkg.in/yaml.v3"
)

type dialogType int

const (
	plainText dialogType = iota
	multiChoice
	textBox
)

const (
	plainTextField   = "plain_text"
	multiChoiceField = "multi_choice"
	textBoxField     = "text_box"

	entryDialogLabel = "entry_dialog"
	lastDialogLabel  = "last_dialog"
	nextDialogLabel  = "next_dialog"
	textLabel        = "text"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "Usage: textgame <game_file>")
		return
	}

	text, err := os.ReadFile(os.Args[1])
	if err != nil {
		panic(err)
	}

	var game map[string]interface{}
	err = yaml.Unmarshal(text, &game)
	if err != nil {
		panic(err)
	}

	dialog, ok := game[entryDialogLabel].(map[string]interface{})
	if !ok {
		fmt.Fprintln(os.Stderr, entryDialogLabel, "key missing")
	}

	reader := bufio.NewReader(os.Stdin)
	nextDialog := entryDialogLabel
	for nextDialog != lastDialogLabel {
		text, ok := dialog["text"].(string)
		if !ok {
			panic(fmt.Sprintln("Dialog requires", textLabel, "field"))
		}
		fmt.Println(text)
		dt := determineDialogType(dialog)
		switch dt {
		case plainText:
			fmt.Println("(Press Enter to continue)")
			reader.ReadString('\n')

			nextDialog = dialog[plainTextField].(map[string]interface{})[nextDialogLabel].(string)
			dialog, ok = game[nextDialog].(map[string]interface{})
			if !ok {
				panic(fmt.Sprintln(os.Stderr, nextDialog, "key missing"))
			}
		case multiChoice:
			options := dialog[multiChoiceField].([]interface{})
			numOptions := len(options)
			if numOptions == 0 {
				panic(fmt.Sprintf("%s: a multi-choice dialog requires at least one option\n", nextDialog))
			}

			for i, option := range options {
				fmt.Printf("%d. %s\n", i+1, option.(map[string]interface{})[textLabel].(string))
			}
			fmt.Printf("(Choose 1...%d)\n", numOptions)

			var input string
			var chosenOption int
			for {
				fmt.Scan(&input)
				chosenOption, err = strconv.Atoi(input)
				if err != nil || chosenOption < 1 || chosenOption > numOptions {
					fmt.Printf("(Choose 1...%d)\n", numOptions)
				} else {
					break
				}
			}
			chosenOption--

			nextDialog = options[chosenOption].(map[string]interface{})[nextDialogLabel].(string)
			dialog, ok = game[nextDialog].(map[string]interface{})
			if !ok {
				panic(fmt.Sprintln(os.Stderr, nextDialog, "key missing"))
			}
		case textBox:
			options := dialog[textBoxField].([]interface{})
			if len(options) == 0 {
				panic(fmt.Sprintf("%s: a text box dialog requires at least one answer\n", nextDialog))
			}

		inputLoop:
			for {
				fmt.Println("(Enter text)")
				var input string
				fmt.Scan(&input)

				for _, option := range options {
					if input != option.(map[string]interface{})[textLabel].(string) {
						continue
					}
					nextDialog = option.(map[string]interface{})[nextDialogLabel].(string)
					dialog, ok = game[nextDialog].(map[string]interface{})
					if !ok {
						panic(fmt.Sprintln(os.Stderr, nextDialog, "key missing"))
					}
					break inputLoop
				}
			}
		}
	}

	lastText, ok := dialog[textLabel].(string)
	if !ok {
		panic(fmt.Sprintln("Dialog requires", textLabel, "field"))
	}
	fmt.Println(lastText)
}

func determineDialogType(dialog map[string]interface{}) dialogType {
	_, hasPlainText := dialog[plainTextField]
	_, hasMultiChoice := dialog[multiChoiceField]
	_, hasTextBox := dialog[textBoxField]
	switch {
	case hasPlainText && !hasMultiChoice && !hasTextBox:
		return plainText
	case !hasPlainText && hasMultiChoice && !hasTextBox:
		return multiChoice
	case !hasPlainText && !hasMultiChoice && hasTextBox:
		return textBox
	default:
		panic(fmt.Sprintf("A dialog requires exactly one of%s, %s or %s\n", plainTextField, multiChoiceField, textBoxField))
	}
}
