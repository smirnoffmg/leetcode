package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"
)

const origin = "https://leetcode.com"

func problemURL(slug string) string { return origin + "/problems/" + slug + "/" }

// client — авторизованный клиент leetcode. Куки берутся из браузера и живут около недели,
// автологина по паролю у leetcode нет.
type client struct {
	session string
	csrf    string
	http    *http.Client
}

func newClient(session, csrf string) *client {
	return &client{session: session, csrf: csrf, http: &http.Client{Timeout: 30 * time.Second}}
}

var errUnauthorized = errors.New("leetcode не принял куки: обнови LEETCODE_SESSION и LEETCODE_CSRF")

func (c *client) do(method, url, referer string, body, out any) error {
	var payload []byte
	if body != nil {
		var err error
		if payload, err = json.Marshal(body); err != nil {
			return err
		}
	}

	req, err := http.NewRequest(method, url, bytes.NewReader(payload))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Origin", origin)
	req.Header.Set("Referer", referer)
	req.Header.Set("User-Agent", "Mozilla/5.0")
	req.Header.Set("X-Csrftoken", c.csrf)
	req.Header.Set("Cookie", fmt.Sprintf("LEETCODE_SESSION=%s; csrftoken=%s", c.session, c.csrf))

	resp, err := c.http.Do(req)
	if err != nil {
		return err
	}
	defer func() { _ = resp.Body.Close() }()

	switch resp.StatusCode {
	case http.StatusOK:
	case http.StatusUnauthorized, http.StatusForbidden:
		return errUnauthorized
	case http.StatusTooManyRequests:
		return errors.New("leetcode просит подождать: слишком много запросов")
	default:
		return fmt.Errorf("leetcode ответил %s на %s", resp.Status, url)
	}

	return json.NewDecoder(resp.Body).Decode(out)
}

const questionIDQuery = `query q($slug: String!) {
  question(titleSlug: $slug) { questionId questionFrontendId title }
}`

type question struct {
	ID         string `json:"questionId"`
	FrontendID string `json:"questionFrontendId"`
	Title      string `json:"title"`
}

// questionID — внутренний id задачи, он не совпадает с номером из URL, а submit хочет именно его.
func (c *client) question(slug string) (question, error) {
	var payload struct {
		Data struct {
			Question *question `json:"question"`
		} `json:"data"`
		Errors []struct {
			Message string `json:"message"`
		} `json:"errors"`
	}
	body := map[string]any{"query": questionIDQuery, "variables": map[string]any{"slug": slug}}
	if err := c.do(http.MethodPost, origin+"/graphql", problemURL(slug), body, &payload); err != nil {
		return question{}, err
	}
	if len(payload.Errors) > 0 {
		return question{}, fmt.Errorf("graphql: %s", payload.Errors[0].Message)
	}
	if payload.Data.Question == nil {
		return question{}, fmt.Errorf("задача %q не найдена", slug)
	}
	return *payload.Data.Question, nil
}

func (c *client) submit(slug, questionID, lang, code string) (int64, error) {
	body := map[string]any{
		"lang":         lang,
		"questionSlug": slug,
		"question_id":  questionID,
		"typed_code":   code,
	}
	var resp struct {
		SubmissionID int64  `json:"submission_id"`
		Error        string `json:"error"`
	}
	url := fmt.Sprintf("%s/problems/%s/submit/", origin, slug)
	if err := c.do(http.MethodPost, url, problemURL(slug), body, &resp); err != nil {
		return 0, err
	}
	if resp.Error != "" {
		return 0, errors.New(resp.Error)
	}
	if resp.SubmissionID == 0 {
		return 0, errors.New("leetcode не вернул submission_id")
	}
	return resp.SubmissionID, nil
}

type checkResult struct {
	State             string      `json:"state"`
	StatusMsg         string      `json:"status_msg"`
	TotalCorrect      int         `json:"total_correct"`
	TotalTestcases    int         `json:"total_testcases"`
	StatusRuntime     string      `json:"status_runtime"`
	StatusMemory      string      `json:"status_memory"`
	RuntimePercentile float64     `json:"runtime_percentile"`
	MemoryPercentile  float64     `json:"memory_percentile"`
	LastTestcase      string      `json:"last_testcase"`
	CodeOutput        looseString `json:"code_output"`
	ExpectedOutput    string      `json:"expected_output"`
	CompileError      string      `json:"compile_error"`
	RuntimeError      string      `json:"runtime_error"`
}

func (r checkResult) accepted() bool { return r.StatusMsg == "Accepted" }

// waitResult поллит вердикт: leetcode отдаёт state PENDING/STARTED, пока считает.
func (c *client) waitResult(slug string, id int64, timeout time.Duration) (checkResult, error) {
	url := fmt.Sprintf("%s/submissions/detail/%d/check/", origin, id)
	deadline := time.Now().Add(timeout)

	for {
		var result checkResult
		if err := c.do(http.MethodGet, url, problemURL(slug), nil, &result); err != nil {
			return checkResult{}, err
		}
		if result.State == "SUCCESS" {
			return result, nil
		}
		if time.Now().After(deadline) {
			return checkResult{}, fmt.Errorf("вердикт не пришёл за %s, состояние %q", timeout, result.State)
		}
		time.Sleep(time.Second)
	}
}

// looseString — code_output приходит то строкой, то массивом строк.
type looseString string

func (s *looseString) UnmarshalJSON(data []byte) error {
	var str string
	if err := json.Unmarshal(data, &str); err == nil {
		*s = looseString(str)
		return nil
	}
	var list []string
	if err := json.Unmarshal(data, &list); err != nil {
		return err
	}
	if len(list) > 0 {
		*s = looseString(list[len(list)-1])
	}
	return nil
}
