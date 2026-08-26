package services

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
)

type FonntePayload struct {
	Target      string `json:"target"`
	Message     string `json:"message"`
	CountryCode string `json:"countryCode"`
}

func SendFonnteWhatsApp(target string, message string) error {
	token := os.Getenv("FONNTE_TOKEN")
	if token == "" {
		log.Println("[FONNTE] WARNING: FONNTE_TOKEN not set in environment. Skipping WhatsApp notification.")
		return fmt.Errorf("FONNTE_TOKEN not configured")
	}

	payload := FonntePayload{
		Target:      target,
		Message:     message,
		CountryCode: "62",
	}

	jsonData, err := json.Marshal(payload)
	if err != nil {
		log.Printf("[FONNTE] Error marshaling payload: %v", err)
		return fmt.Errorf("failed to marshal payload: %w", err)
	}

	req, err := http.NewRequest("POST", "https://api.fonnte.com/send", bytes.NewBuffer(jsonData))
	if err != nil {
		log.Printf("[FONNTE] Error creating request: %v", err)
		return fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Authorization", token)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Printf("[FONNTE] Error sending request: %v", err)
		return fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Printf("[FONNTE] API returned status %d for target %s", resp.StatusCode, target)
		return fmt.Errorf("fonnte API returned status %d", resp.StatusCode)
	}

	log.Printf("[FONNTE] ✅ WhatsApp notification sent successfully to %s", target)
	return nil
}

func BuildDonationMessage(campaignTitle string, amount float64) string {
	return fmt.Sprintf(
		"Terimakasih telah mendonasi, bantuan anda sangat membantu dan berharga bagi mereka.\n\n"+
			"Untuk menyelesaikan proses donasi pada campaign *%s*, mohon lakukan transfer sebesar *Rp%s* ke rekening bersama kami:\n"+
			"Bank: BCA\n"+
			"No. Rekening: 123456789\n"+
			"a/n OpenDonasi Escrow\n\n"+
			"Setelah transfer, harap unggah bukti transfer di web kami.",
		campaignTitle,
		formatRupiahGo(amount),
	)
}


func formatRupiahGo(amount float64) string {
	intAmount := int64(amount)
	str := fmt.Sprintf("%d", intAmount)

	
	n := len(str)
	if n <= 3 {
		return str
	}

	var result []byte
	for i, c := range str {
		if i > 0 && (n-i)%3 == 0 {
			result = append(result, '.')
		}
		result = append(result, byte(c))
	}
	return string(result)
}
