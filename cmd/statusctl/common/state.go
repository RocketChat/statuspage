package common

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"os/user"
	"path"
	"strings"

	"github.com/RocketChat/statuscentral/client/oauthclient"
)

type State struct {
	Session *oauthclient.ClientSession `json:"session"`
}

var stateFile = "$HOME/.statusctl/config"

func getStateFile() (string, error) {
	usr, err := user.Current()
	if err != nil {
		return "", err
	}

	return strings.ReplaceAll(stateFile, "$HOME", usr.HomeDir), nil
}

func LoadState() (*State, error) {
	filePath, err := getStateFile()
	if err != nil {
		return nil, err
	}

	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return nil, err
	}

	jsonFile, err := ioutil.ReadFile(filePath)
	if err != nil {
		log.Printf("LoadState read err   #%v ", err)
		return nil, err
	}

	state := &State{}

	if err := json.Unmarshal(jsonFile, state); err != nil {
		log.Fatalf("Unmarshal: %v", err)
		return nil, err
	}

	return state, nil
}

func SaveState(state State) error {
	filePath, err := getStateFile()
	if err != nil {
		return err
	}

	jsonState, err := json.Marshal(state)
	if err != nil {
		log.Fatalf("Unmarshal: %v", err)
		return err
	}

	directory, _ := path.Split(filePath)

	if err := os.MkdirAll(directory, 0700); err != nil { //nolint:gomnd // Tech debt
		return err
	}

	if err := ioutil.WriteFile(filePath, jsonState, 0700); err != nil { //nolint:gomnd,gosec // Need to check why we need this perm if at all
		return err
	}

	return nil
}

func DeleteState() error {
	filePath, err := getStateFile()
	if err != nil {
		return err
	}

	if err := os.Remove(filePath); err != nil {
		return err
	}

	return nil
}
