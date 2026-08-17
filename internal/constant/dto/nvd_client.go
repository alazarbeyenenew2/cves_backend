package dto

import (
	"time"

	"github.com/alazarbeyenenew2/cves_backend/platform/logger"
)

type NVDClient struct {
	ApiKey     string
	BaseURL    string
	PageSize   int
	MaxRetries int
	Logger     logger.Logger
}

// Root Response
type NVDResponse struct {
	TotalResults    int             `json:"totalResults"`
	StartIndex      int             `json:"startIndex"`
	Version         string          `json:"version"`
	Vulnerabilities []Vulnerability `json:"vulnerabilities"`
}

// Vulnerability
type Vulnerability struct {
	CVE CVE `json:"cve"`
}

// CVE
type CVE struct {
	ID             string          `json:"id"`
	Published      string          `json:"published"`
	LastModified   string          `json:"lastModified"`
	Descriptions   []Description   `json:"descriptions"`
	Metrics        Metrics         `json:"metrics"`
	Weaknesses     []Weakness      `json:"weaknesses"`
	Configurations []Configuration `json:"configurations"`
	References     []Reference     `json:"references"`
	Affected       any             `json:"affected"`
}

// Description
type Description struct {
	Lang  string `json:"lang"`
	Value string `json:"value"`
}

// Metrics
type Metrics struct {
	CvssMetricV31 []CVSSMetricV31 `json:"cvssMetricV31"`
	CvssMetricV30 []CVSSMetricV30 `json:"cvssMetricV30"`
	CvssMetricV2  []CVSSMetricV2  `json:"cvssMetricV2"`
}

type CVSSMetricV31 struct {
	CvssData CVSSDataV3 `json:"cvssData"`
}

type CVSSMetricV30 struct {
	CvssData CVSSDataV3 `json:"cvssData"`
}
type CVSSDataV3 struct {
	BaseScore    float64 `json:"baseScore"`
	BaseSeverity string  `json:"baseSeverity"`
	AttackVector string  `json:"attackVector"`
	VectorString string  `json:"vectorString"`
}

type CVSSMetricV2 struct {
	BaseSeverity string     `json:"baseSeverity"`
	CvssData     CVSSDataV2 `json:"cvssData"`
}

type CVSSDataV2 struct {
	BaseScore    float64 `json:"baseScore"`
	AccessVector string  `json:"accessVector"`
	VectorString string  `json:"vectorString"`
}

type Weakness struct {
	Description []WeaknessDescription `json:"description"`
}

type WeaknessDescription struct {
	Value string `json:"value"`
}

type Configuration struct {
	Nodes []Node `json:"nodes"`
}

type Node struct {
	CpeMatch []CPEMatch `json:"cpeMatch"`
}

type CPEMatch struct {
	Criteria string `json:"criteria"`
}

type Reference struct {
	URL string `json:"url"`
}
type Meta struct {
	LastFetchTime time.Time `json:"lastFetchTime"`
	TotalFetched  int       `json:"totalFetched"`
	SyncStatus    string    `json:"syncStatus"`
}
