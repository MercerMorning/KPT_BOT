package clients

import (
	"KPT_BOT/config"
	"context"
	"encoding/json"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/option"
	"google.golang.org/api/sheets/v4"
	"log"
	"os"
	"strconv"
)

type SheetsClient struct {
	Code   string
	Bot    *tgbotapi.BotAPI
	Update tgbotapi.Update
	ChatId int64
}

type GoogleOAuthCredentials struct {
	Installed struct {
		ClientID                string   `json:"client_id"`
		ProjectID               string   `json:"project_id"`
		AuthURI                 string   `json:"auth_uri"`
		TokenURI                string   `json:"token_uri"`
		AuthProviderX509CertURL string   `json:"auth_provider_x509_cert_url"`
		ClientSecret            string   `json:"client_secret"`
		RedirectURIs            []string `json:"redirect_uris"`
	} `json:"installed"`
}

func (gc *SheetsClient) getConfig() *oauth2.Config {
	b, err := os.ReadFile("credentials.json")
	if err != nil {
		log.Fatalf("Unable to read client secret file: %v", err)
	}
	creds := GoogleOAuthCredentials{}
	json.Unmarshal(b, &creds)
	creds.Installed.RedirectURIs = []string{
		"http://" + config.Config("HOST") + ":" + config.Config("PORT") + "/" + strconv.FormatInt(gc.ChatId, 10),
	}
	b, err = json.Marshal(creds)
	if err != nil {
		panic(err)
	}
	// If modifying these scopes, delete your previously saved token.json.
	config, err := google.ConfigFromJSON(b, "https://www.googleapis.com/auth/spreadsheets")
	if err != nil {
		log.Fatalf("Unable to parse client secret file to config: %v", err)
	}
	return config
}

func (gc *SheetsClient) getTokenFromWeb(config *oauth2.Config) {
	authURL := config.AuthCodeURL("state-token", oauth2.AccessTypeOffline)
	msg := tgbotapi.NewMessage(gc.Update.Message.Chat.ID,
		fmt.Sprintf("Go to the following link in your browser then type the "+
			"authorization code: \n%v\n", authURL),
	)
	if _, err := gc.Bot.Send(msg); err != nil {
		panic(err)
	}
}

func (gc *SheetsClient) getTokenFromWebWithCode(config *oauth2.Config) *oauth2.Token {
	authCode := gc.Code
	tok, err := config.Exchange(context.TODO(), authCode)
	if err != nil {
		fmt.Printf("Unable to retrieve token from web: %v", err)
	}
	return tok
}

func (gc *SheetsClient) RequestCode() {
	gc.getTokenFromWeb(gc.getConfig())
}

func (gc *SheetsClient) GetToken() *oauth2.Token {
	return gc.getTokenFromWebWithCode(gc.getConfig())

}

func (gc *SheetsClient) InitTable(tok *oauth2.Token, spreadsheetId string) {
	// If modifying these scopes, delete your previously saved token.json.
	client := gc.getConfig().Client(context.Background(), tok)

	ctx := context.Background()
	srv, err := sheets.NewService(ctx, option.WithHTTPClient(client))
	if err != nil {
		log.Fatalf("Unable to retrieve Sheets client: %v", err)
	}

	rangeData := "sheet1!A1:F1"
	values := [][]interface{}{{"Ситуация", "Мысль", "Эмоция", "Ощущение", "Действие", "Дата"}}
	rb := &sheets.BatchUpdateValuesRequest{
		ValueInputOption: "USER_ENTERED",
	}
	rb.Data = append(rb.Data, &sheets.ValueRange{
		Range:  rangeData,
		Values: values,
	})

	_, err = srv.Spreadsheets.Values.BatchUpdate(
		spreadsheetId, rb,
	).Do()
	if err != nil {
		log.Fatalf("Unable to retrieve data from sheet: %v", err)
	}
}

func (gc *SheetsClient) Append(tok *oauth2.Token, spreadsheetId string, data []string) {
	client := gc.getConfig().Client(context.Background(), tok)

	ctx := context.Background()
	srv, err := sheets.NewService(ctx, option.WithHTTPClient(client))

	if err != nil {
		log.Fatalf("Unable to retrieve Sheets client: %v", err)
	}

	//readRange := "Class Data!A2:E"
	var records [][]interface{}

	// Слайс строк, который хотим добавить
	//stringSlice := []string{
	//	usrSession.Diary.Situation,
	//	usrSession.Diary.Thought,
	//	usrSession.Diary.Emotion,
	//	usrSession.Diary.Feeling,
	//	usrSession.Diary.Action,
	//	time.Now().Format(time.DateOnly),
	//}

	// Конвертируем []string в []interface{}
	interfaceSlice := make([]interface{}, len(data))
	for i, v := range data {
		interfaceSlice[i] = v
	}

	// Добавляем в records
	records = append(records, interfaceSlice)

	valueInputOption := "USER_ENTERED"
	insertDataOption := "INSERT_ROWS"
	rb := &sheets.ValueRange{
		Values: records,
	}

	response, err := srv.Spreadsheets.Values.Append(spreadsheetId, "A:F", rb).ValueInputOption(valueInputOption).InsertDataOption(insertDataOption).Context(ctx).Do()
	if err != nil || response.HTTPStatusCode != 200 {
		panic(err)
	}
}
