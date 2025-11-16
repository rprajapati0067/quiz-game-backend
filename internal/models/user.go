package models

type User struct {
    ID        string `dynamodbav:"user_id"`
    Name      string `dynamodbav:"name"`
    Phone     string `dynamodbav:"phone"`
    Email     string `dynamodbav:"email"`
    Verified  bool   `dynamodbav:"verified"`
    Blocked   bool   `dynamodbav:"blocked"`
    Points    int64  `dynamodbav:"points"`
}
