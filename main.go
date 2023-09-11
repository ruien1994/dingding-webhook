package main

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"strings"
)

type Alert struct {
	Text string `json:"text"`
}

type DingTalkMessage struct {
	MsgType string            `json:"msgtype"`
	Text    DingTalkText      `json:"text"`
	At      DingTalkAt        `json:"at"`
}

type DingTalkText struct {
	Content string `json:"content"`
}

type DingTalkAt struct {
	AtAll bool `json:"isAtAll"`
}

func main() {
	http.HandleFunc("/webhook/", webhookHandler)
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func webhookHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	// 解析出 ACCESS_TOKEN
	path := strings.TrimPrefix(r.URL.Path, "/webhook/")
	access_token := strings.TrimPrefix(path, "/")

	var alert Alert
	err := json.NewDecoder(r.Body).Decode(&alert)
	if err != nil {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	log.Printf("Received alert: %s\n", alert.Text)

	// 转发告警信息给钉钉机器人
	sendDingTalkAlert(alert.Text, access_token)

	w.WriteHeader(http.StatusOK)
}

func sendDingTalkAlert(text string, access_token string) {
	message := DingTalkMessage{
		MsgType: "text",
		Text: DingTalkText{
			Content: text,
		},
		At: DingTalkAt{
			AtAll: true, // @所有人
		},
	}

	payload, err := json.Marshal(message)
	if err != nil {
		log.Println("Failed to marshal DingTalk message:", err)
		return
	}

	// 构造钉钉机器人的 Webhook URL
	webhookURL := "https://oapi.dingtalk.com/robot/send?access_token=" + access_token

	resp, err := http.Post(webhookURL, "application/json", bytes.NewBuffer(payload))
	if err != nil {
		log.Println("Failed to send DingTalk alert:", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Println("DingTalk alert failed with status code:", resp.StatusCode)
		return
	}

	log.Println("DingTalk alert sent successfully")
}

