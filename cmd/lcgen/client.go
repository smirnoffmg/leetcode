package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

const api = "https://leetcode.com/graphql"

type Question struct {
	FrontendID       string `json:"questionFrontendId"`
	Title            string `json:"title"`
	TitleSlug        string `json:"titleSlug"`
	Difficulty       string `json:"difficulty"`
	Content          string `json:"content"`
	Hints            []string
	MetaData         string `json:"metaData"`
	ExampleTestcases string `json:"exampleTestcases"`
	SimilarQuestions string `json:"similarQuestions"`
	Stats            string `json:"stats"`
	TopicTags        []struct {
		Name string `json:"name"`
	} `json:"topicTags"`
	CodeSnippets []struct {
		LangSlug string `json:"langSlug"`
		Code     string `json:"code"`
	} `json:"codeSnippets"`
}

func (q Question) goSnippet() string {
	for _, s := range q.CodeSnippets {
		if s.LangSlug == "golang" {
			return s.Code
		}
	}
	return "func solve() {\n}\n"
}

func (q Question) acceptance() string {
	var stats struct {
		AcRate string `json:"acRate"`
	}
	_ = json.Unmarshal([]byte(q.Stats), &stats)
	return stats.AcRate
}

type SimilarQuestion struct {
	Title      string `json:"title"`
	TitleSlug  string `json:"titleSlug"`
	Difficulty string `json:"difficulty"`
}

func (q Question) similar() []SimilarQuestion {
	var out []SimilarQuestion
	_ = json.Unmarshal([]byte(q.SimilarQuestions), &out)
	return out
}

type Meta struct {
	Name   string `json:"name"`
	Params []struct {
		Name string `json:"name"`
		Type string `json:"type"`
	} `json:"params"`
	Return *struct {
		Type string `json:"type"`
	} `json:"return"`
}

func (q Question) meta() Meta {
	var m Meta
	_ = json.Unmarshal([]byte(q.MetaData), &m)
	return m
}

const dailyQuery = `query { activeDailyCodingChallengeQuestion { date question { titleSlug } } }`

const questionQuery = `query q($slug: String!) {
  question(titleSlug: $slug) {
    questionFrontendId title titleSlug difficulty content hints metaData
    exampleTestcases similarQuestions stats
    topicTags { name }
    codeSnippets { langSlug code }
  }
}`

func graphQL(query string, vars map[string]any, out any) error {
	body, err := json.Marshal(map[string]any{"query": query, "variables": vars})
	if err != nil {
		return err
	}

	req, err := http.NewRequest(http.MethodPost, api, bytes.NewReader(body))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	// Без Referer/User-Agent leetcode иногда отвечает 403.
	req.Header.Set("Referer", "https://leetcode.com/")
	req.Header.Set("User-Agent", "Mozilla/5.0")

	resp, err := (&http.Client{Timeout: 30 * time.Second}).Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("leetcode ответил %s", resp.Status)
	}

	var payload struct {
		Data   json.RawMessage `json:"data"`
		Errors []struct {
			Message string `json:"message"`
		} `json:"errors"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&payload); err != nil {
		return err
	}
	if len(payload.Errors) > 0 {
		return fmt.Errorf("graphql: %s", payload.Errors[0].Message)
	}
	return json.Unmarshal(payload.Data, out)
}

func fetchDaily() (date, slug string, err error) {
	var data struct {
		Daily struct {
			Date     string `json:"date"`
			Question struct {
				TitleSlug string `json:"titleSlug"`
			} `json:"question"`
		} `json:"activeDailyCodingChallengeQuestion"`
	}
	if err := graphQL(dailyQuery, nil, &data); err != nil {
		return "", "", err
	}
	return data.Daily.Date, data.Daily.Question.TitleSlug, nil
}

func fetchQuestion(slug string) (Question, error) {
	var data struct {
		Question *Question `json:"question"`
	}
	if err := graphQL(questionQuery, map[string]any{"slug": slug}, &data); err != nil {
		return Question{}, err
	}
	if data.Question == nil {
		return Question{}, fmt.Errorf("задача %q не найдена", slug)
	}
	return *data.Question, nil
}
