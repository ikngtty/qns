package qns

import (
	"fmt"
	"net/http"
	"os"

	"github.com/PuerkitoBio/goquery"
)

type notification struct {
	Kind string
	Text string
}

func Load() {
	req, err := http.NewRequest("GET", "https://qiita.com/notifications", nil)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	req.AddCookie(&http.Cookie{
		Name:  "user_session_key",
		Value: os.Getenv("qiita_user_session_key")})
	req.AddCookie(&http.Cookie{
		Name:  "secure_token",
		Value: os.Getenv("qiita_secure_token")})

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer resp.Body.Close()

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	notifications := make([]notification, 0)
	doc.Find(".notification_actionWrapper").Each(func(index int, wrapper *goquery.Selection) {
		contents := selectionToSlice(wrapper.Contents())
		var kind string
		if len(contents) == 7 && contents[4].HasClass("bold") && contents[4].Text() == "LGTM" {
			kind = "LGTM"
		} else if len(contents) == 7 && contents[4].HasClass("bold") && contents[4].Text() == "ストック" {
			kind = "ストック"
		} else if len(contents) == 7 && contents[4].HasClass("bold") && contents[4].Text() == "コメント" {
			kind = "コメント"
		} else if len(contents) == 7 && contents[4].HasClass("bold") && contents[4].Text() == "編集リクエスト" {
			kind = "編集リクエスト"
		} else if len(contents) == 7 && contents[4].HasClass("bold") && contents[4].Text() == "リンク" {
			kind = "リンク"
		} else if len(contents) == 7 && contents[4].HasClass("bold") && contents[4].Text() == "編集" {
			kind = "編集"
		} else {
			kind = "その他"
		}
		notifications = append(notifications, notification{kind, wrapper.Text()})
	})

	for _, notifi := range notifications {
		fmt.Printf("kind:%s text:%s\n", notifi.Kind, notifi.Text)
	}
}

func selectionToSlice(sel *goquery.Selection) []*goquery.Selection {
	sels := make([]*goquery.Selection, sel.Length())
	sel.Each(func(index int, sel *goquery.Selection) {
		sels[index] = sel
	})
	return sels
}
