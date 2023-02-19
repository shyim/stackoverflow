package main

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v4"
	"log"
	"os"
	"strconv"
	"strings"
	"time"
)

func main() {
	ctx := context.Background()
	conn, err := pgx.Connect(ctx, "postgres://stackoverflow:stackoverflow@localhost:5432/stackoverflow")

	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	}

	defer conn.Close(ctx)

	updateQuestions(ctx, conn, "shopware")
	updateQuestions(ctx, conn, "shopware6")
}

func updateQuestions(ctx context.Context, conn *pgx.Conn, tag string) {
	page := 1
	for {
		log.Printf("Fetching questions for tag: %s in page: %d\n", tag, page)
		questions, err := GetQuestions(ctx, page, tag)

		if err != nil {
			log.Fatal(err)
		}

		ids := make([]string, 0)

		for _, question := range questions.Items {
			updateAuthor(conn, ctx, question.Owner)
			updateQuestion(conn, ctx, question)

			ids = append(ids, strconv.FormatInt(question.QuestionId, 10))
		}

		updateAnswers(conn, ctx, ids)

		log.Println("Waiting some time to not get rate limited")
		time.Sleep(time.Second * 30)

		page++

		if !questions.HasMore {
			break
		}

		// one page is currently enough
		break
	}
}

func updateAnswers(conn *pgx.Conn, ctx context.Context, ids []string) {
	page := 1

	for {
		answers, err := GetAnswersOfQuestion(ctx, strings.Join(ids, ";"), page)

		if err != nil {
			log.Fatalln(err)
		}

		sql := "insert into answer (id, question_id, is_accepted, score, last_activity_date, creation_date, content_license, body, creator_id)\nvalues ($1, $2, $3, $4, $5, $6, $7, $8, $9)\nON CONFLICT (id)\nDO UPDATE\nSET\n    is_accepted = excluded.is_accepted,\n    score = excluded.score,\n    last_activity_date = excluded.last_activity_date,\n    body = excluded.body\n       "

		for _, answer := range answers.Items {
			updateAuthor(conn, ctx, answer.Owner)

			_, err := conn.Exec(
				ctx,
				sql,
				answer.AnswerId,
				answer.QuestionId,
				answer.IsAccepted,
				answer.Score,
				formatDate(answer.LastActivityDate),
				formatDate(answer.CreationDate),
				answer.ContentLicense,
				answer.Body,
				answer.Owner.AccountId,
			)

			if err != nil {
				log.Fatalln(err)
			}
		}

		page++

		if !answers.HasMore {
			break
		}
	}
}

func updateAuthor(conn *pgx.Conn, ctx context.Context, owner StackoverflowOwner) {
	sql := "insert into \"user\" (id, reputation, user_id, user_type, profile_image, display_name, link)\nvalues ($1, $2, $3, $4, $5, $6, $7)\nON CONFLICT (id)\nDO UPDATE\n       SET reputation = excluded.reputation,\n       user_id = excluded.user_id,\n       user_type = excluded.user_type,\n       profile_image = excluded.profile_image,\n       display_name = excluded.display_name"

	_, err := conn.Exec(
		ctx,
		sql,
		owner.AccountId,
		owner.Reputation,
		owner.UserId,
		owner.UserType,
		owner.ProfileImage,
		owner.DisplayName,
		owner.Link,
	)

	if err != nil {
		log.Fatalln(err)
	}
}

func updateQuestion(conn *pgx.Conn, ctx context.Context, question StackoverflowListingElement) {
	questionSQL := "INSERT INTO question (id, creator_id, is_answered, view_count, closed_date, score, last_activity_date,\n                      creation_date, last_edit_date, link, closed_reason, title, content_license, accepted_answer_id, body)\nvalues ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15)\nON CONFLICT (id)\nDO UPDATE \n    SET is_answered = excluded.is_answered,\n       view_count = excluded.view_count,\n       closed_date = excluded.closed_date,\n       score = excluded.score,\n       last_activity_date = excluded.last_activity_date,\n       closed_reason = excluded.closed_reason,\n       accepted_answer_id = excluded.accepted_answer_id,\n       body = excluded.body\n       "

	_, err := conn.Exec(
		ctx,
		questionSQL,
		question.QuestionId,
		question.Owner.AccountId,
		question.IsAnswered,
		question.ViewCount,
		formatDate(question.ClosedDate),
		question.Score,
		formatDate(&question.LastActivityDate),
		formatDate(&question.CreationDate),
		formatDate(question.LastEditDate),
		question.Link,
		question.ClosedReason,
		question.Title,
		question.ContentLicense,
		question.AcceptedAnswerId,
		question.Body,
	)

	if err != nil {
		log.Fatal(err)
	}

	_, _ = conn.Exec(ctx, "DELETE FROM question_to_tag WHERE question_id = $1", question.QuestionId)

	for _, tag := range question.Tags {
		// Yolo insert. We don't care if it exists
		conn.Exec(ctx, "insert into tag (tag) values ($1)", tag)

		var tagId int64
		conn.QueryRow(ctx, "SELECT id FROM tag WHERE tag = $1", tag).Scan(&tagId)

		conn.Exec(ctx, "INSERT INTO question_to_tag (question_id, tag_id) VALUES ($1, $2)", question.QuestionId, tagId)
	}
}

func formatDate(timestamp *int64) *time.Time {
	if timestamp == nil {
		return nil
	}

	uff := time.Unix(*timestamp, 0)

	return &uff
}
