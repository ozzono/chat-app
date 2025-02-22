package bot

import (
	"chat-app/internal/models"
	"encoding/csv"
	"fmt"
	"net/http"
	"strings"
	"time"
)

const (
	helpMenu = "These are the available commands:\n/help: shows this help menu\n/stock: fetches the value of a given stock using the format '/stock=APPL.US'"
)

func ProcessCMD(input string) (models.Message, error) {

	msg := models.Message{
		Timestamp: time.Now().UTC(),
		Nickname:  "BOT",
	}

	switch {
	case strings.HasPrefix(input, "/help"):
		msg.Content = helpMenu
		return msg, nil
	case strings.HasPrefix(input, "/stock"):
		msg.Content = getStock(parseCMD(input))
		return msg, nil
	default:
		msg.Content = "invalid command;\n" + helpMenu
		return msg, nil
	}
}

func parseCMD(input string) string {
	return strings.TrimPrefix(strings.Split(input, " ")[0], "/stock=")
}

func getStock(code string) string {
	code = strings.ToUpper(code)
	res, err := http.Get(fmt.Sprintf("https://stooq.com/q/l/?s=%s&f=sd2t2ohlcv&h&e=csv", code))
	if err != nil {
		return fmt.Sprintf("failed to fetch %s stock value; please try again.", code)
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return fmt.Sprintf("failed to fetch %s stock value; please try again.", code)
	}

	reader := csv.NewReader(res.Body)
	records, err := reader.ReadAll()
	if err != nil {
		return fmt.Sprintf("failed to fetch %s stock value; please contact system admin.", code)
	}

	if len(records) < 2 {
		return fmt.Sprintf("failed to fetch %s stock value; please contact system admin.", code)
	}

	data := records[1]
	price := data[6]
	if price == "N/D" {
		return fmt.Sprintf("%s stock not found; do you want to try another one?", code)
	}

	return fmt.Sprintf("%s value is $%s per unit", code, price)
}
