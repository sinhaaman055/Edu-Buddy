package models

import (
	"time"
)

type Room struct {
	ID                 string               `json:"id" bson:"_id"`
	RoomName           string               `json:"roomname" bson:"roomname"`
	Description        string               `json:"description" bson:"description"`
	CreatedBy          string               `json:"createdby" bson:"createdby"`
	CreatedAt          time.Time            `json:"createdat" bson:"createdat"`
	Members            []string             `json:"members" bson:"members"`
	Password           string               `json:"password" bson:"password"` 
}
type TestSession struct{
	TestId   string 
	TestIsActive bool
	TimeRemaiming int
    TickerDone  chan bool
	Questions  []MockQuestion
	Response    map[string]map[string]string
}
type MockQuestion struct{
    QuestionId    string       `json:"questionId" bson:"question_id"`
	Text          string       `json:"text" bson:"text"`
	Options     []string       `json:"options" bson:"options"`
	Answer        string       `json:"answer" bson:"answer"`
}
type TestHistory struct {
    ID          string               `json:"id" bson:"_id"`
    RoomID      string               `json:"roomId" bson:"room_id"`
    TestID      string               `json:"testId" bson:"test_id"`
    ExamName    string               `json:"examName" bson:"exam_name"`
    TotalTime   int                  `json:"totalTime" bson:"total_time"` 
    CreatedAt   time.Time            `json:"createdAt" bson:"created_at"`
    Questions   []MockQuestion       `json:"questions" bson:"questions"`
    UserPerformances map[string]UserStats `json:"userPerformances" bson:"user_performances"`
}
type UserStats struct {
    Username       string            `json:"username" bson:"username"`
    Score          int               `json:"score" bson:"score"`
    CorrectCount   int               `json:"correctCount" bson:"correct_count"`
    WrongCount     int               `json:"wrongCount" bson:"wrong_count"`
    TickedAnswers  map[string]string `json:"tickedAnswers" bson:"ticked_answers"`
}