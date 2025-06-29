package session

type Session struct {
	Stage   int
	Code    string
	SheetId string
	Diary   Diary
}

type Diary struct {
	Situation string
	Thought   string
	Emotion   string
	Feeling   string
	Action    string
}

const (
	Start                     = 0
	WritingDiary              = 1
	RequestGoogleSheetsApiKey = 2
	ReceiveGoogleSheetsApiKey = 3
	WritingDiarySituation     = 4
	WritingDiaryThought       = 5
	WritingDiaryEmotion       = 6
	WritingDiaryFeeling       = 7
	WritingDiaryAction        = 8
)
