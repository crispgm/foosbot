// Package app .
package app

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-lark/lark"
	larkgin "github.com/go-lark/lark-gin"

	"github.com/crispgm/foosbot/internal/def"
)

// CardValue .
type CardValue struct {
	Action string
}

// LoadRoutes .
func LoadRoutes(r *gin.Engine) {
	bot := newBot()
	bot.StartHeartbeat()

	mw := larkgin.NewLarkMiddleware()

	g := r.Group("/lark")
	{
		g.Use(mw.LarkChallengeHandler())

		g.GET("/ping", func(c *gin.Context) {
			c.String(http.StatusOK, "pong")
		})

		eventGroup := g.Group("/event")
		{
			eventGroup.Use(mw.LarkEventHandler())
			mw.WithTokenVerification(def.AppVerificationToken)
			eventGroup.POST("/callback", func(c *gin.Context) {
				if event, ok := mw.GetEvent(c); ok {
					switch event.Header.EventType {
					case lark.EventTypeMessageReceived:
						if msg, err := event.GetMessageReceived(); err == nil {
							if msg.Message.MessageType == lark.MsgText {
								var content lark.TextContent
								_ = json.Unmarshal([]byte(msg.Message.Content), &content)
								log.Println(msg.Sender.SenderID.OpenID, "sended:", content.Text)

								if msg.Sender.SenderID.OpenID == def.AdminOpenID {
									if content.Text == "notify" || content.Text == "1" {
										notifyPlayers(bot, LevelNormal)
									} else if content.Text == "notify more" {
										notifyPlayers(bot, LevelExtended)
									}
								}
							}
						} else {
							log.Println(err)
						}

					case lark.EventTypeMessageReactionCreated:
						if evt, err := event.GetMessageReactionCreated(); err == nil {
							msgResp, err := bot.WithUserIDType(lark.UIDOpenID).GetMessage(evt.MessageID)
							if err != nil {
								log.Println(err)
								break
							}
							if msgResp.Data.Items[0].Sender.ID != def.AppID {
								break
							}
							log.Println("Create reaction:", evt.ReactionType.EmojiType)
							if evt.ReactionType.EmojiType == string(lark.EmojiTypeOK) || evt.ReactionType.EmojiType == string(lark.EmojiTypeJIAYI) {
								_ = replyToAction(bot, evt.UserID.OpenID, evt.MessageID, "+1")
							} else if evt.ReactionType.EmojiType == string(lark.EmojiTypeMinusOne) {
								_ = replyToAction(bot, evt.UserID.OpenID, evt.MessageID, "-1")
							}
						} else {
							log.Println(err)
						}
					default:
						// just ignore
					}
				}
			})
		}
		cardGroup := g.Group("/card")
		{
			cardGroup.Use(mw.LarkCardHandler())
			cardGroup.POST("/callback", func(c *gin.Context) {
				if card, ok := mw.GetCardCallback(c); ok {
					action := card.Action
					var value CardValue
					_ = json.Unmarshal([]byte(action.Value), &value)
					log.Println("Received:", action.Tag, action.Option, value.Action)
					if action.Tag == "button" {
						err := replyToAction(bot, card.OpenID, card.MessageID, value.Action)
						if err != nil {
							log.Println(err)
						}
					} else if action.Tag == "select_person" {
						openID := action.Option
						resp, err := bot.GetUserInfo(lark.WithOpenID(openID))
						if err != nil {
							log.Println(err)
							return
						}
						if value.Action == "buzz" {
							err = notifySingle(bot, resp.Data.User.EnterpriseEmail, openID)
							if err != nil {
								log.Println(err)
							}
						} else if value.Action == "buzz_phone" {
							bot.WithUserIDType(lark.UIDOpenID)
							resp, err := bot.BuzzMessage(lark.BuzzTypePhone, card.MessageID, openID)
							if err != nil {
								log.Println(err)
							}
							if resp.Code != 0 {
								log.Println(resp.Code, resp.Msg)
							}
						}
					}
				}
			})
		}
	}
}
