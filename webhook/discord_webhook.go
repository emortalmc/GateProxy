package webhook

import (
	"bytes"
	"fmt"
	"go.minekube.com/gate/pkg/edition/java/proxy"
	"net/http"
)

func SendWebhookMessage(payload []byte, url string) {
	if url == "" {
		return
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(payload))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("User-Agent", "Go-Discord")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return
	}

	defer resp.Body.Close()
}

func PlayerJoined(plr proxy.Player, plrCount int, url string) {
	playersText := "players"
	if plrCount == 1 {
		playersText = "player"
	}

	var jsonData = []byte(fmt.Sprintf(`{
		"username": "%s",
		"content": "Joined the server! (%d %s)",
		"avatar_url": "https://mc-heads.net/avatar/%s/100"
	}`, plr.Username(), plrCount, playersText, plr.ID().String()))

	go SendWebhookMessage(jsonData, url)
}

func PlayerLeft(plr proxy.Player, plrCount int, url string) {
	playersText := "players"
	if plrCount == 1 {
		playersText = "player"
	}

	var jsonData = []byte(fmt.Sprintf(`{
		"username": "%s",
		"content": "Left the server! (%d %s)",
		"avatar_url": "https://mc-heads.net/avatar/%s/100"
	}`, plr.Username(), plrCount, playersText, plr.ID().String()))

	go SendWebhookMessage(jsonData, url)
}
