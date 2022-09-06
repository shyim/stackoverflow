package main

type StackoverflowListingCollection struct {
	Items          []StackoverflowListingElement `json:"items"`
	HasMore        bool                          `json:"has_more"`
	QuotaMax       int                           `json:"quota_max"`
	QuotaRemaining int                           `json:"quota_remaining"`
}

type StackoverflowListingElement struct {
	Tags             []string           `json:"tags"`
	Owner            StackoverflowOwner `json:"owner"`
	IsAnswered       bool               `json:"is_answered"`
	ViewCount        int                `json:"view_count"`
	ClosedDate       *int64             `json:"closed_date,omitempty"`
	AnswerCount      int                `json:"answer_count"`
	Score            int                `json:"score"`
	LastActivityDate int64              `json:"last_activity_date"`
	CreationDate     int64              `json:"creation_date"`
	LastEditDate     *int64             `json:"last_edit_date,omitempty"`
	QuestionId       int64              `json:"question_id"`
	Link             string             `json:"link"`
	ClosedReason     string             `json:"closed_reason,omitempty"`
	Title            string             `json:"title"`
	ContentLicense   string             `json:"content_license,omitempty"`
	AcceptedAnswerId *int               `json:"accepted_answer_id,omitempty"`
	Body             string             `json:"body_markdown"`
}

type StackoverflowOwner struct {
	AccountId    int    `json:"account_id"`
	Reputation   int    `json:"reputation"`
	UserId       int    `json:"user_id"`
	UserType     string `json:"user_type"`
	ProfileImage string `json:"profile_image"`
	DisplayName  string `json:"display_name"`
	Link         string `json:"link"`
}

type StackoverflowAnswerCollection struct {
	Items          []StackoverflowAnswerElement `json:"items"`
	HasMore        bool                         `json:"has_more"`
	QuotaMax       int                          `json:"quota_max"`
	QuotaRemaining int                          `json:"quota_remaining"`
}

type StackoverflowAnswerElement struct {
	Owner            StackoverflowOwner `json:"owner"`
	IsAccepted       bool               `json:"is_accepted"`
	Score            int                `json:"score"`
	LastActivityDate *int64             `json:"last_activity_date"`
	CreationDate     *int64             `json:"creation_date"`
	AnswerId         int                `json:"answer_id"`
	QuestionId       int                `json:"question_id"`
	ContentLicense   string             `json:"content_license"`
	Body             string             `json:"body_markdown"`
}
