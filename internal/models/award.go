package models

type Award struct {
    ID        string `dynamodbav:"award_id"`
    Product   string `dynamodbav:"product"`
    PointCost int64  `dynamodbav:"point_cost"`
}
