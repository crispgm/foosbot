// Package def .
package def

import (
	"errors"
	"fmt"
	"log"
	"net/mail"
	"os"
	"strings"

	"github.com/joho/godotenv"
)

// App credentials
var (
	Port                 string
	AppID                string
	AppSecret            string
	AppVerificationToken string
	ChatID               string
	AdminOpenID          string
	NotifyNormalUsers    []string
	NotifyExtendedUsers  []string
)

// LoadVariables load vars from dotenv
func LoadVariables() error {
	mode := os.Getenv("FOOSBOT_MODE")
	if mode != "production" {
		mode = "testing"
	}
	log.Println("Mode is", mode)

	path := fmt.Sprintf(".env.%s", mode)
	err := godotenv.Load(path)
	if err != nil {
		log.Printf(".env.%s not found. Use real environment variables.", mode)
	}

	Port = os.Getenv("FOOSBOT_PORT")
	AppID = os.Getenv("FOOSBOT_APP_ID")
	AppSecret = os.Getenv("FOOSBOT_APP_SECRET")
	AppVerificationToken = os.Getenv("FOOSBOT_APP_VERIFICATION_TOKEN")
	ChatID = os.Getenv("FOOSBOT_CHAT_ID")
	AdminOpenID = os.Getenv("FOOSBOT_ADMIN_OPEN_ID")

	notifyNormalUsers := os.Getenv("FOOSBOT_NOTIFY_NORMAL_USERS")
	NotifyNormalUsers = strings.Split(notifyNormalUsers, ",")
	err = validUsers(NotifyNormalUsers...)
	if err != nil {
		return err
	}
	notifyExtendedUsers := os.Getenv("FOOSBOT_NOTIFY_EXTENDED_USERS")
	NotifyExtendedUsers = strings.Split(notifyExtendedUsers, ",")
	err = validUsers(NotifyExtendedUsers...)
	if err != nil {
		return err
	}

	log.Printf("Loaded variables: mode=%s, port=%s, app_id=%s, app_secret=%s, app_verification_token=%s, chat_id=%s, admin_open_id=%s, notify{normal_users=%s extended_users=%s}\n", mode, Port, AppID, AppSecret, AppVerificationToken, ChatID, AdminOpenID, notifyNormalUsers, notifyExtendedUsers)
	return nil
}

func validUsers(users ...string) error {
	if len(users) == 0 {
		return errors.New("Empty user group")
	}

	for _, email := range users {
		_, err := mail.ParseAddress(email)
		if err != nil {
			return err
		}
	}

	return nil
}
