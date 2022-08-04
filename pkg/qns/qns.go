package qns

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path"
	"time"

	"github.com/PuerkitoBio/goquery"
)

type notification struct {
	Kind string
	Text string
}

type LoadSettings struct {
	Pages int
}

func Load(settings LoadSettings) {
	err := os.MkdirAll(getQnsDirPath(), 0777)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	notifications := make([]notification, 0)
	for page := 1; page <= settings.Pages; page++ {
		url := fmt.Sprintf("https://qiita.com/notifications?page=%d", page)
		req, err := http.NewRequest("GET", url, nil)
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

		found := false
		doc.Find(".notification_actionWrapper").Each(func(index int, wrapper *goquery.Selection) {
			found = true

			contents := selectionToSlice(wrapper.Contents())
			var kind string
			if len(contents) == 5 && contents[2].HasClass("bold") && contents[2].Text() == "フォロー" {
				kind = "フォロー"
			} else if len(contents) == 5 && contents[2].HasClass("bold") && contents[2].Text() == "採用" {
				kind = "採用"
			} else if len(contents) == 7 && contents[4].HasClass("bold") && contents[4].Text() == "LGTM" {
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

		if !found {
			break
		}

		time.Sleep(time.Second)
	}

	fmt.Printf("loaded %d notifications!\n", len(notifications))

	notificationsJson, err := json.Marshal(notifications)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	err = os.WriteFile(getNotificationsPath(), notificationsJson, 0664)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func View(kind string) {
	notificationsJson, err := os.ReadFile(getNotificationsPath())
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	var notifications []notification
	err = json.Unmarshal(notificationsJson, &notifications)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	for _, notifi := range notifications {
		if notifi.Kind != kind {
			continue
		}
		fmt.Println(notifi.Text)
	}
}

func getQnsDirPath() string {
	home, err := os.UserHomeDir()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	return path.Join(home, ".qns")
}

func getNotificationsPath() string {
	return path.Join(getQnsDirPath(), "notifications.json")
}

func selectionToSlice(sel *goquery.Selection) []*goquery.Selection {
	sels := make([]*goquery.Selection, sel.Length())
	sel.Each(func(index int, sel *goquery.Selection) {
		sels[index] = sel
	})
	return sels
}
