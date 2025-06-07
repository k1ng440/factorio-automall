package recipe

import (
    "encoding/json"
    "os"
	"io"
)
type Recipe struct {
		AllowProductivity bool   `json:"allow_productivity"`
		Category          string `json:"category"`
		Ingredients       []struct {
			Amount float64    `json:"amount"`
			Name   string `json:"name"`
		} `json:"ingredients"`
		Key           string `json:"key"`
		LocalizedName struct {
			En string `json:"en"`
		} `json:"localized_name"`
		Order   string `json:"order"`
		Results []struct {
			Amount float64    `json:"amount"`
			Name   string `json:"name"`
		} `json:"results"`
		Subgroup string `json:"subgroup"`
	}
type FactorioData struct {
	Items []struct {
		Group         string `json:"group"`
		IconCol       int    `json:"icon_col"`
		IconRow       int    `json:"icon_row"`
		Key           string `json:"key"`
		LocalizedName struct {
			En string `json:"en"`
		} `json:"localized_name"`
		Order     string `json:"order"`
		Subgroup  string `json:"subgroup"`
		Type      string `json:"type"`
		StackSize int    `json:"stack_size,omitempty"`
	} `json:"items"`
	Recipes []*Recipe `json:"recipes"`
}


func ReadFile (fileName string) (*FactorioData, error) {
    file, err := os.Open(fileName)
    if err != nil {
        return nil, err
    }
    defer file.Close()

    data, err := io.ReadAll(file)
    if err != nil {
        return nil, err
    }

    var factorioData FactorioData
    err = json.Unmarshal(data, &factorioData)
    if err != nil {
        return nil, err
    }

    return &factorioData, nil
}

