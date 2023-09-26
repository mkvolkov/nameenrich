package enrich

import (
	"context"
	"io"
	"nameenrich/types"
	"net/http"
	"testing"
	"time"

	"github.com/goccy/go-json"
	"github.com/stretchr/testify/assert"
)

func TestAgeEnrich(t *testing.T) {
	assert := assert.New(t)

	var tests = []struct {
		name        string
		expectedAge int
	}{
		{"Andrei", 40},
		{"Alina", 39},
		{"Feodor", 52},
	}

	for _, test := range tests {
		httpClient := http.Client{}

		ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
		defer cancel()

		url := "https://api.agify.io/?name=" + test.name

		req, _ := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)

		resp, _ := httpClient.Do(req)

		data, _ := io.ReadAll(resp.Body)

		var agedata types.AgeData = types.AgeData{}

		json.Unmarshal(data, &agedata)

		assert.Equal(agedata.Age, test.expectedAge)
	}
}

func TestGenderEnrich(t *testing.T) {
	assert := assert.New(t)

	var tests = []struct {
		name           string
		expectedGender string
	}{
		{"Andrei", "male"},
		{"Alina", "female"},
		{"Feodor", "male"},
	}

	for _, test := range tests {
		httpClient := http.Client{}

		ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
		defer cancel()

		url := "https://api.genderize.io/?name=" + test.name

		req, _ := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)

		resp, _ := httpClient.Do(req)

		data, _ := io.ReadAll(resp.Body)

		var gendata types.GenderData = types.GenderData{}

		json.Unmarshal(data, &gendata)

		assert.Equal(gendata.Gender, test.expectedGender)
	}
}

func TestNationEnrich(t *testing.T) {
	assert := assert.New(t)

	var name string = "Andrei"
	var Countries []types.Nationality = []types.Nationality{
		{
			Country:     "RO",
			Probability: 0.514,
		},
		{
			Country:     "MD",
			Probability: 0.12,
		},
		{
			Country:     "BY",
			Probability: 0.092,
		},
		{
			Country:     "RU",
			Probability: 0.05,
		},
		{
			Country:     "UA",
			Probability: 0.028,
		},
	}

	httpClient := http.Client{}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	url := "https://api.nationalize.io/?name=" + name

	req, _ := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)

	resp, _ := httpClient.Do(req)

	data, _ := io.ReadAll(resp.Body)

	var nationData types.NationData = types.NationData{}

	json.Unmarshal(data, &nationData)

	assert.Equal(Countries, nationData.Nationalities)
}
