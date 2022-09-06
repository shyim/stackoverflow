package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

func GetQuestions(ctx context.Context, page int, tag string) (*StackoverflowListingCollection, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, fmt.Sprintf("https://api.stackexchange.com/2.3/questions?page=%d&order=desc&sort=activity&tagged=%s&site=stackoverflow&filter=!nKzQUR30W7&pagesize=100", page, tag), nil)

	if err != nil {
		return nil, err
	}

	resp, err := http.DefaultClient.Do(req)

	if err != nil {
		return nil, err
	}

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("rate limited")
	}

	var collection StackoverflowListingCollection

	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		return nil, err
	}

	if err := json.Unmarshal(body, &collection); err != nil {
		return nil, err
	}

	return &collection, nil
}

func GetAnswersOfQuestion(ctx context.Context, questionIds string, page int) (*StackoverflowAnswerCollection, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, fmt.Sprintf("https://api.stackexchange.com/2.3/questions/%s/answers?page=%d&order=desc&sort=activity&site=stackoverflow&filter=!)qRpaqDpVQcGRJQynhhv&pagesize=100", questionIds, page), nil)

	if err != nil {
		return nil, err
	}

	resp, err := http.DefaultClient.Do(req)

	if err != nil {
		return nil, err
	}

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("rate limited")
	}

	var collection StackoverflowAnswerCollection

	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		return nil, err
	}

	if err := json.Unmarshal(body, &collection); err != nil {
		return nil, err
	}

	return &collection, nil
}
