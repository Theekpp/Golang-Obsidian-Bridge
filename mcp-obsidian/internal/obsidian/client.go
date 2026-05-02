package obsidian

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

// Client interacts with Obsidian's Local REST API plugin.
type Client struct {
	baseURL    string
	apiKey     string
	httpClient *http.Client
}

// NewClient creates a new Obsidian REST API client.
// baseURL example: "http://127.0.0.1:27123"
// apiKey is the token from the Local REST API plugin settings.
func NewClient(baseURL, apiKey string) *Client {
	return &Client{
		baseURL: strings.TrimRight(baseURL, "/"),
		apiKey:  apiKey,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

func (c *Client) do(method, path string, body any) ([]byte, int, error) {
	var bodyReader io.Reader
	if body != nil {
		data, err := json.Marshal(body)
		if err != nil {
			return nil, 0, fmt.Errorf("marshal request body: %w", err)
		}
		bodyReader = bytes.NewReader(data)
	}

	req, err := http.NewRequest(method, c.baseURL+path, bodyReader)
	if err != nil {
		return nil, 0, fmt.Errorf("create request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+c.apiKey)
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, 0, fmt.Errorf("execute request: %w", err)
	}
	defer resp.Body.Close()

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, resp.StatusCode, fmt.Errorf("read response: %w", err)
	}

	return data, resp.StatusCode, nil
}

func (c *Client) doText(method, path, contentType, textBody string) ([]byte, int, error) {
	var bodyReader io.Reader
	if textBody != "" {
		bodyReader = strings.NewReader(textBody)
	}

	req, err := http.NewRequest(method, c.baseURL+path, bodyReader)
	if err != nil {
		return nil, 0, fmt.Errorf("create request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+c.apiKey)
	if contentType != "" {
		req.Header.Set("Content-Type", contentType)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, 0, fmt.Errorf("execute request: %w", err)
	}
	defer resp.Body.Close()

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, resp.StatusCode, fmt.Errorf("read response: %w", err)
	}

	return data, resp.StatusCode, nil
}

// ---- Types ----------------------------------------------------------------

// FileEntry represents a file in the vault listing.
type FileEntry struct {
	Path    string `json:"path"`
	Type    string `json:"type"` // "file" | "dir"
	Content string `json:"content,omitempty"`
}

// VaultListResponse is returned by GET /vault/
type VaultListResponse struct {
	Files []string `json:"files"`
}

// SearchResult represents a single note match from the search API.
type SearchResult struct {
	Filename string        `json:"filename"`
	Score    float64       `json:"score"`
	Matches  []SearchMatch `json:"matches,omitempty"`
}

// SearchMatch is an individual text match inside a note.
type SearchMatch struct {
	Match   SearchMatchRange `json:"match"`
	Context string           `json:"context"`
}

// SearchMatchRange holds start/end character offsets of a match.
type SearchMatchRange struct {
	Start int `json:"start"`
	End   int `json:"end"`
}

// ActiveNoteResponse is returned by GET /active/
type ActiveNoteResponse struct {
	Path    string `json:"path"`
	Content string `json:"content"`
	Tags    []string `json:"tags"`
}

// NoteResponse is returned by GET /vault/{path}
type NoteResponse struct {
	Path    string   `json:"path"`
	Content string   `json:"content"`
	Tags    []string `json:"tags"`
	Stat    NoteStat `json:"stat"`
}

// NoteStat holds file metadata.
type NoteStat struct {
	CTime int64 `json:"ctime"`
	MTime int64 `json:"mtime"`
	Size  int64 `json:"size"`
}

// PeriodicNoteType defines the type of periodic note.
type PeriodicNoteType string

const (
	PeriodDaily   PeriodicNoteType = "daily"
	PeriodWeekly  PeriodicNoteType = "weekly"
	PeriodMonthly PeriodicNoteType = "monthly"
	PeriodYearly  PeriodicNoteType = "yearly"
)

// ---- API Methods ----------------------------------------------------------

// ListFiles returns all file paths in the vault (or a subdirectory).
func (c *Client) ListFiles(dirPath string) ([]string, error) {
	apiPath := "/vault/"
	if dirPath != "" {
		apiPath += url.PathEscape(strings.TrimPrefix(dirPath, "/"))
		if !strings.HasSuffix(apiPath, "/") {
			apiPath += "/"
		}
	}

	data, status, err := c.do(http.MethodGet, apiPath, nil)
	if err != nil {
		return nil, err
	}
	if status != http.StatusOK {
		return nil, fmt.Errorf("list files: HTTP %d: %s", status, string(data))
	}

	var resp VaultListResponse
	if err := json.Unmarshal(data, &resp); err != nil {
		return nil, fmt.Errorf("parse list response: %w", err)
	}

	return resp.Files, nil
}

// ReadNote returns the content of a note by its vault path.
func (c *Client) ReadNote(notePath string) (*NoteResponse, error) {
	apiPath := "/vault/" + url.PathEscape(strings.TrimPrefix(notePath, "/"))
	data, status, err := c.doText(http.MethodGet, apiPath, "", "")
	if err != nil {
		return nil, err
	}
	if status == http.StatusNotFound {
		return nil, fmt.Errorf("note not found: %s", notePath)
	}
	if status != http.StatusOK {
		return nil, fmt.Errorf("read note: HTTP %d: %s", status, string(data))
	}

	// Try JSON first (newer plugin versions), fall back to raw text.
	var note NoteResponse
	if err := json.Unmarshal(data, &note); err != nil {
		note.Path = notePath
		note.Content = string(data)
	}

	return &note, nil
}

// CreateOrUpdateNote creates or fully replaces a note.
func (c *Client) CreateOrUpdateNote(notePath, content string) error {
	apiPath := "/vault/" + url.PathEscape(strings.TrimPrefix(notePath, "/"))
	_, status, err := c.doText(http.MethodPut, apiPath, "text/markdown", content)
	if err != nil {
		return err
	}
	if status != http.StatusOK && status != http.StatusNoContent && status != http.StatusCreated {
		return fmt.Errorf("create/update note: HTTP %d", status)
	}

	return nil
}

// AppendToNote appends text to an existing note.
func (c *Client) AppendToNote(notePath, content string) error {
	apiPath := "/vault/" + url.PathEscape(strings.TrimPrefix(notePath, "/"))
	_, status, err := c.doText(http.MethodPost, apiPath, "text/markdown", content)
	if err != nil {
		return err
	}
	if status != http.StatusOK && status != http.StatusNoContent && status != http.StatusCreated {
		return fmt.Errorf("append to note: HTTP %d", status)
	}

	return nil
}

// DeleteNote removes a note from the vault.
func (c *Client) DeleteNote(notePath string) error {
	apiPath := "/vault/" + url.PathEscape(strings.TrimPrefix(notePath, "/"))
	_, status, err := c.do(http.MethodDelete, apiPath, nil)
	if err != nil {
		return err
	}
	if status == http.StatusNotFound {
		return fmt.Errorf("note not found: %s", notePath)
	}
	if status != http.StatusOK && status != http.StatusNoContent {
		return fmt.Errorf("delete note: HTTP %d", status)
	}

	return nil
}

// SearchNotes performs a simple full-text search across the vault.
func (c *Client) SearchNotes(query string) ([]SearchResult, error) {
	apiPath := "/search/simple/?query=" + url.QueryEscape(query)
	data, status, err := c.do(http.MethodPost, apiPath, nil)
	if err != nil {
		return nil, err
	}
	if status != http.StatusOK {
		return nil, fmt.Errorf("search: HTTP %d: %s", status, string(data))
	}

	var results []SearchResult
	if err := json.Unmarshal(data, &results); err != nil {
		return nil, fmt.Errorf("parse search response: %w", err)
	}

	return results, nil
}

// GetActiveNote returns the note currently open in Obsidian.
func (c *Client) GetActiveNote() (*ActiveNoteResponse, error) {
	data, status, err := c.do(http.MethodGet, "/active/", nil)
	if err != nil {
		return nil, err
	}
	if status == http.StatusNotFound {
		return nil, fmt.Errorf("no active note open in Obsidian")
	}
	if status != http.StatusOK {
		return nil, fmt.Errorf("get active note: HTTP %d: %s", status, string(data))
	}

	var note ActiveNoteResponse
	if err := json.Unmarshal(data, &note); err != nil {
		return nil, fmt.Errorf("parse active note: %w", err)
	}

	return &note, nil
}

// OpenNote opens a note in the Obsidian UI.
func (c *Client) OpenNote(notePath string) error {
	apiPath := "/open/" + url.PathEscape(strings.TrimPrefix(notePath, "/"))
	_, status, err := c.do(http.MethodPost, apiPath, nil)
	if err != nil {
		return err
	}
	if status != http.StatusOK && status != http.StatusNoContent {
		return fmt.Errorf("open note: HTTP %d", status)
	}

	return nil
}

// GetPeriodicNote returns a periodic note (daily/weekly/monthly/yearly).
func (c *Client) GetPeriodicNote(period PeriodicNoteType) (*ActiveNoteResponse, error) {
	apiPath := fmt.Sprintf("/periodic/%s/", period)
	data, status, err := c.do(http.MethodGet, apiPath, nil)
	if err != nil {
		return nil, err
	}
	if status == http.StatusNotFound {
		return nil, fmt.Errorf("no %s periodic note found", period)
	}
	if status != http.StatusOK {
		return nil, fmt.Errorf("get periodic note: HTTP %d: %s", status, string(data))
	}

	var note ActiveNoteResponse
	if err := json.Unmarshal(data, &note); err != nil {
		return nil, fmt.Errorf("parse periodic note: %w", err)
	}

	return &note, nil
}

// ServerInfo returns information about the Obsidian REST API server.
func (c *Client) ServerInfo() (map[string]any, error) {
	data, status, err := c.do(http.MethodGet, "/", nil)
	if err != nil {
		return nil, err
	}
	if status != http.StatusOK {
		return nil, fmt.Errorf("server info: HTTP %d: %s", status, string(data))
	}

	var info map[string]any
	if err := json.Unmarshal(data, &info); err != nil {
		return nil, fmt.Errorf("parse server info: %w", err)
	}

	return info, nil
}
