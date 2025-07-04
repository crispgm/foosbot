package app

import (
	"context"
	"errors"
	"fmt"
	"log"
	"math/rand"
	"strings"

	"github.com/go-lark/lark/v2"

	"github.com/crispgm/foosbot/internal/def"
)

// notify levels
const (
	LevelNormal = iota + 1
	LevelExtended
)

func buildCard(chatID string, users ...string) lark.OutcomingMessage {
	var at []string
	for _, person := range users {
		at = append(at, fmt.Sprintf("<at email=\"%s\"></at>", person))
	}
	b := lark.NewCardBuilder()
	card := b.Card(
		b.Div(
			b.Field(b.Text(chickenSoup())),
		),
		b.Action(
			b.Button(b.Text("+1")).Value(map[string]interface{}{"action": "+1"}).Primary(),
			b.Button(b.Text("+1，3分钟后到")).Value(map[string]interface{}{"action": "+1，3分钟后到"}).Primary(),
			b.Button(b.Text("+1，5分钟后到")).Value(map[string]interface{}{"action": "+1，5分钟后到"}).Primary(),
			b.Button(b.Text("+1，10分钟后到")).Value(map[string]interface{}{"action": "+1，10分钟后到"}).Primary(),
		),
		b.Action(
			b.Button(b.Text("-1")).Value(map[string]interface{}{"action": "-1，我菜，所以缺席训练"}).Danger().Confirm("确认", "我菜，所以缺席训练"),
			b.Button(b.Text("-1，出差中")).Value(map[string]interface{}{"action": "-1，出差中，建议弯道"}),
		),
		b.Hr(),
		b.Action(
			b.Button(b.Text("签到")).Value(map[string]interface{}{"action": "已到达现场"}),
			b.Button(b.Text("速来 1=3")).Value(map[string]interface{}{"action": "我已到现场，请速来训练！"}),
			b.Button(b.Text("速来 3=1")).Value(map[string]interface{}{"action": "速来！9999=1"}),
		),
		b.Hr(),
		b.Div(
			b.Field(b.Text("催一下")),
		),
		b.Action(
			b.SelectMenu().SelectPerson().Value(map[string]interface{}{"action": "buzz"}),
		),
		b.Hr(),
		b.Div(
			b.Field(b.Text(strings.Join(at, " ")).LarkMd()),
		),
	).
		Title("1?")
	msg := lark.
		NewMsgBuffer(lark.MsgInteractive).
		BindChatID(chatID).
		Card(card.String()).
		Build()
	return msg
}

func notifySingle(ctx context.Context, bot *lark.Bot, email, openID string) error {
	b := lark.NewCardBuilder()
	card := b.Card(
		b.Div(
			b.Field(b.Text(fmt.Sprintf("<at email=\"%s\"></at>", email)).LarkMd()),
		),
		b.Hr(),
		b.Action(
			b.Button(b.Text("+1")).Value(map[string]interface{}{"action": "+1"}).Primary(),
			b.Button(b.Text("+1，3分钟后到")).Value(map[string]interface{}{"action": "+1，3分钟后到"}).Primary(),
			b.Button(b.Text("+1，5分钟后到")).Value(map[string]interface{}{"action": "+1，5分钟后到"}).Primary(),
			b.Button(b.Text("+1，10分钟后到")).Value(map[string]interface{}{"action": "+1，10分钟后到"}).Primary(),
		),
		b.Action(
			b.Button(b.Text("-1")).Value(map[string]interface{}{"action": "-1，我菜，所以缺席训练"}).Danger().Confirm("确认", "我菜，所以缺席训练"),
			b.Button(b.Text("-1，出差中")).Value(map[string]interface{}{"action": "-1，出差中，建议弯道"}),
		),
		b.Hr(),
		b.Action(
			b.Button(b.Text("签到")).Value(map[string]interface{}{"action": "已到达现场"}),
		),
		b.Hr(),
		b.Div(b.Field(b.Text("电话加急"))),
		b.Action(
			b.SelectMenu().SelectPerson().Value(map[string]interface{}{"action": "buzz_phone"}).Confirm("确认", "确认进行电话加急"),
		),
	).
		Title("速来训练！")
	msg := lark.
		NewMsgBuffer(lark.MsgInteractive).
		BindChatID(def.ChatID).
		Card(card.String()).
		Build()
	resp, err := bot.PostMessage(ctx, msg)
	if err != nil {
		return err
	}
	if resp.Code != 0 {
		log.Println(resp.Code, resp.Msg)
		return errors.New(resp.Msg)
	}

	// buzz
	bot.WithUserIDType(lark.UIDOpenID)
	buzzResp, err := bot.BuzzMessage(ctx, lark.BuzzTypeInApp, resp.Data.MessageID, openID)
	if err != nil {
		log.Println(err)
		return err
	}
	if buzzResp.Code != 0 {
		log.Println(buzzResp.Code, buzzResp.Msg)
		return errors.New(resp.Msg)
	}

	return nil
}

// notifyPlayers send notification to players
func notifyPlayers(ctx context.Context, bot *lark.Bot, level int) error {
	var users []string
	if level == LevelNormal {
		// Notify (normal)
		users = def.NotifyNormalUsers
	} else if level == LevelExtended {
		// Notify extended
		users = def.NotifyExtendedUsers
	} else {
		return errors.New("no users given")
	}

	// do send
	msg := buildCard(def.ChatID, users...)
	resp, err := bot.PostMessage(ctx, msg)
	if err != nil {
		return err
	}
	if resp.Code != 0 {
		log.Println(resp.Code, resp.Msg)
		return errors.New(resp.Msg)
	}

	return nil
}

// replyToAction replies to user who reacts
func replyToAction(ctx context.Context, bot *lark.Bot, openID, msgID, action string) error {
	userResp, err := bot.GetUserInfo(ctx, lark.WithOpenID(openID))
	if err != nil {
		return err
	}
	if userResp.Code != 0 {
		log.Println(userResp.Code, userResp.Msg)
		return errors.New(userResp.Msg)
	}
	name := userResp.Data.User.Name

	msgText := fmt.Sprintf("%s: %s", name, action)
	if strings.HasPrefix(action, "+1") || strings.HasPrefix(action, "签到") {
		msgText = fmt.Sprintf("%s: %s\n随机序号: %d", name, action, rand.Int()%100)
	}
	msg := lark.
		NewMsgBuffer(lark.MsgText).
		Text(msgText).
		BindOpenID(openID).
		BindReply(msgID).
		Build()
	msgResp, err := bot.ReplyMessage(ctx, msg)
	if err != nil {
		return err
	}
	if msgResp.Code != 0 {
		log.Println(msgResp.Code, msgResp.Msg)
		return errors.New(msgResp.Msg)
	}

	return nil
}

// Buzz!
