package common

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
)

// StringPrompt asks for a string value using the label
func StringPrompt(label string) string {
	var s string
	r := bufio.NewReader(os.Stdin)
	for {
		fmt.Fprint(os.Stderr, label+" ")
		s, _ = r.ReadString('\n')
		if s != "" {
			break
		}
	}
	return strings.TrimSpace(s)
}

// StringPromptWithDefault asks for a string value using the label if none uses default value
func StringPromptWithDefault(label string, defaultValue string) string {
	result := StringPrompt(label)

	if result == "" {
		return defaultValue
	}

	return result
}

// IntPrompt asks for a Int value using the label
func IntPrompt(label string, defaultValue int) (int, error) {
	result := StringPrompt(label)

	if result == "" {
		return defaultValue, nil
	}

	num, err := strconv.Atoi(result)
	if err != nil {
		return -1, err
	}

	return num, nil
}

func GetYesNoPrompt(label string, yes bool) (bool, error) {
	yesText := "y"
	noText := "n"
	defaultSelection := noText

	if yes {
		yesText = "Y"
		defaultSelection = yesText
	} else {
		noText = "N"
		defaultSelection = noText
	}

	yesNoSelection := StringPromptWithDefault(fmt.Sprintf("%s [%s/%s]:", label, yesText, noText), defaultSelection)

	YesNo := strings.ToLower(yesNoSelection)

	if YesNo != "y" && YesNo != "yes" && YesNo != "n" && YesNo != "no" {
		return false, errors.New("invalid selection")
	}

	if YesNo == "n" || YesNo == "no" {
		return false, nil
	}

	return true, nil
}
