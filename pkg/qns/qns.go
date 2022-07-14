package qns

import (
	"fmt"
	"net/http"
	"os"

	"github.com/PuerkitoBio/goquery"
)

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

	fmt.Println(doc.Html())
}
