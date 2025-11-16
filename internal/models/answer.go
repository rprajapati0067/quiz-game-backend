package models

import "time"

type Answer struct {
    UserID      string    `dynamodbav:"user_id"`
    QuestionID  string    `dynamodbav:"question_id"`
    SubmittedAt time.Time `dynamodbav:"submitted_at"`
    Correct     bool      `dynamodbav:"correct"`
}
