package models

// AssetType defines different asset categories
type AssetType string

const (
	ChartType    AssetType = "chart"
	InsightType  AssetType = "insight"
	AudienceType AssetType = "audience"
)

// Base struct for all assets
type Asset struct {
	ID          string    `json:"id"`
	Type        AssetType `json:"type"`
	Description string    `json:"description"`
}

// Chart-specific struct
type Chart struct {
	Asset
	Title      string   `json:"title"`
	AxesTitles []string `json:"axes_titles"`
	Data       []int    `json:"data"`
}

// Insight-specific struct
type Insight struct {
	Asset
}

// Audience-specific struct
type Audience struct {
	Asset
	Gender       string `json:"gender"`
	BirthCountry string `json:"birth_country"`
	AgeGroup     string `json:"age_group"`
	HoursOnline  int    `json:"hours_online"`
	Purchases    int    `json:"purchases"`
}
