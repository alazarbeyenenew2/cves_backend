package nvd_client

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/alazarbeyenenew2/cves_backend/internal/constant/dto"
	"github.com/alazarbeyenenew2/cves_backend/platform/logger"
	"github.com/spf13/viper"
)

type NVDClient struct {
	apiKey     string
	httpClient *http.Client
	BaseURL    string
	PageSize   int
	MaxRetries int
	logger     logger.Logger
}

func NewNVDClient(req dto.NVDClient) *NVDClient {
	return &NVDClient{
		apiKey: req.ApiKey,
		httpClient: &http.Client{
			Timeout: viper.GetDuration("nvd.timeout"),
		},
		logger:     req.Logger,
		BaseURL:    req.BaseURL,
		PageSize:   req.PageSize,
		MaxRetries: req.MaxRetries,
	}
}

func (c *NVDClient) FetchRange(startDate, endDate time.Time, progressFn func(fetched, total int)) ([]dto.NVDResponse, error) {
	const maxChunkDays = 90

	delay := 6500 * time.Millisecond
	if c.apiKey != "" {
		delay = 700 * time.Millisecond
	}

	var allCVEs []dto.NVDResponse

	currentStart := startDate
	for currentStart.Before(endDate) {
		currentEnd := currentStart.AddDate(0, 0, maxChunkDays)
		if currentEnd.After(endDate) {
			currentEnd = endDate
		}

		startStr := currentStart.UTC().Format("2006-01-02T15:04:05.000")
		endStr := currentEnd.UTC().Format("2006-01-02T15:04:05.000")

		startIndex := 0
		for {
			resp, err := c.fetchPage(startStr, endStr, startIndex)
			if err != nil {
				return nil, err
			}

			allCVEs = append(allCVEs, *resp)

			if progressFn != nil {
				progressFn(len(allCVEs), resp.TotalResults)
			}

			startIndex += c.PageSize
			if startIndex >= resp.TotalResults {
				break
			}

			time.Sleep(delay)
		}

		if currentEnd.Equal(endDate) {
			break
		}
		currentStart = currentEnd.Add(time.Second)
		time.Sleep(delay)
	}

	return allCVEs, nil
}

func (c *NVDClient) fetchPage(startDate, endDate string, startIndex int) (*dto.NVDResponse, error) {
	var lastErr error
	for attempt := 0; attempt < c.MaxRetries; attempt++ {
		if attempt > 0 {
			wait := time.Duration(10*(1<<uint(attempt-1))) * time.Second
			if wait > 60*time.Second {
				wait = 60 * time.Second
			}
			c.logger.Named("nvd_client").Warn(context.Background(), fmt.Sprintf("[NVD] Rate limit encountered (%d/%d): waiting %v before retry...", attempt, c.MaxRetries, wait))
			time.Sleep(wait)
		}
		result, err := c.doFetchPage(startDate, endDate, startIndex)
		if err == nil {
			return result, nil
		}
		lastErr = err

		if isRateLimitOrServerError(err) {
			continue
		}

		return nil, err
	}
	return nil, fmt.Errorf("after %d retries: %w", c.MaxRetries, lastErr)
}

func (c *NVDClient) doFetchPage(startDate, endDate string, startIndex int) (*dto.NVDResponse, error) {
	req, err := http.NewRequest("GET", c.BaseURL, nil)
	if err != nil {
		return nil, err
	}

	q := req.URL.Query()
	q.Set("pubStartDate", startDate)
	q.Set("pubEndDate", endDate)
	q.Set("startIndex", fmt.Sprintf("%d", startIndex))
	q.Set("resultsPerPage", fmt.Sprintf("%d", c.PageSize))
	req.URL.RawQuery = q.Encode()

	if c.apiKey != "" {
		req.Header.Set("apiKey", c.apiKey)
	}
	req.Header.Set("User-Agent", "CVE-Notifier/1.0")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("NVD request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == 429 || resp.StatusCode == 403 {
		return nil, fmt.Errorf("NVD API rate limited (HTTP %d)", resp.StatusCode)
	}
	if resp.StatusCode >= 500 {
		return nil, fmt.Errorf("NVD API server error (HTTP %d)", resp.StatusCode)
	}
	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("NVD API returned HTTP %d", resp.StatusCode)
	}

	var result dto.NVDResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("decode NVD response: %w", err)
	}
	return &result, nil
}

func isRateLimitOrServerError(err error) bool {
	if err == nil {
		return false
	}
	msg := err.Error()
	return strings.Contains(msg, "429") || strings.Contains(msg, "403") || strings.Contains(msg, "503") || strings.Contains(msg, "rate limit") || strings.Contains(msg, "500") || strings.Contains(msg, "502") || strings.Contains(msg, "504")
}

func classifyVulnType(vector, desc string) string {
	v := strings.ToUpper(vector)
	d := strings.ToLower(desc)

	if strings.Contains(d, "sql injection") || strings.Contains(d, "sqli") {
		return "SQL Injection"
	}
	if strings.Contains(d, "cross-site scripting") || strings.Contains(d, "xss") {
		return "XSS"
	}
	if strings.Contains(d, "buffer overflow") || strings.Contains(d, "stack overflow") {
		return "Buffer Overflow"
	}
	if strings.Contains(d, "remote code execution") || strings.Contains(d, "rce") {
		return "RCE"
	}
	if strings.Contains(d, "denial of service") || strings.Contains(d, " dos ") {
		return "DoS"
	}
	if strings.Contains(d, "privilege escalation") || strings.Contains(d, "escalation of privilege") {
		return "Privilege Escalation"
	}
	if strings.Contains(d, "path traversal") || strings.Contains(d, "directory traversal") {
		return "Path Traversal"
	}
	if strings.Contains(d, "csrf") || strings.Contains(d, "cross-site request forgery") {
		return "CSRF"
	}
	if strings.Contains(d, "information disclosure") || strings.Contains(d, "information exposure") {
		return "Info Disclosure"
	}
	if strings.Contains(d, "use after free") {
		return "Use-After-Free"
	}
	if strings.Contains(v, "AV:N") || strings.Contains(v, "NETWORK") {
		return "Remote"
	}
	if strings.Contains(v, "AV:L") || strings.Contains(v, "LOCAL") {
		return "Local"
	}
	return "Other"
}

func (c *NVDClient) FetchSince(since time.Time, progressFn func(fetched, total int)) ([]dto.NVDResponse, error) {
	return c.FetchRange(since, time.Now(), progressFn)
}
