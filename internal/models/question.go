package models

type Question struct {
    ID           string   `dynamodbav:"question_id"`
    Text         string   `dynamodbav:"text"`
    Options      []string `dynamodbav:"options"`
    CorrectIndex int32    `dynamodbav:"correct_index"`
    Slot         int32    `dynamodbav:"slot"`
    CreatedBy    string   `dynamodbav:"created_by"`
}
