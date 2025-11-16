package models

import "time"

type Claim struct {
    UserID   string    `dynamodbav:"user_id"`
    AwardID  string    `dynamodbav:"award_id"`
    Points   int64     `dynamodbav:"points"`
    ClaimedAt time.Time `dynamodbav:"claimed_at"`
}
