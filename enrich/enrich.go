package enrich

import (
	"context"
	"fmt"
	"io"
	"nameenrich/types"
	"net/http"
	"sync"
	"time"

	"github.com/goccy/go-json"
)

func Enrichment(baseMsg *types.MsgBase, enrMsg *types.MsgEnriched) {
	enrMsg.Name = baseMsg.Name
	enrMsg.Surname = baseMsg.Surname
	enrMsg.Patronymic = baseMsg.Patronymic

	urlAge := "https://api.agify.io/?name=" + baseMsg.Name
	urlGender := "https://api.genderize.io/?name=" + baseMsg.Name
	urlNation := "https://api.nationalize.io/?name=" + baseMsg.Name

	var wg sync.WaitGroup

	wg.Add(1)
	go AgeEnrichment(urlAge, &wg, enrMsg)

	wg.Add(1)
	go GenderEnrichment(urlGender, &wg, enrMsg)

	wg.Add(1)
	go NationEnrichment(urlNation, &wg, enrMsg)

	wg.Wait()
}

func AgeEnrichment(url string, wg *sync.WaitGroup, enData *types.MsgEnriched) {
	defer wg.Done()

	httpClient := http.Client{}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		fmt.Println("Error creating request: ", err.Error())
	}

	resp, err := httpClient.Do(req)
	if err != nil {
		fmt.Println("Error doing request: ", err.Error())
	}

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading response body: ", err.Error())
	}

	var agedata types.AgeData = types.AgeData{}

	err = json.Unmarshal(data, &agedata)
	if err != nil {
		fmt.Println("Error unmarshalling age data: ", err.Error())
	}

	enData.Age = agedata.Age
}

func GenderEnrichment(url string, wg *sync.WaitGroup, enData *types.MsgEnriched) {
	defer wg.Done()
	httpClient := http.Client{}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		fmt.Println("Error creating request: ", err.Error())
	}

	resp, err := httpClient.Do(req)
	if err != nil {
		fmt.Println("Error doing request: ", err.Error())
	}

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading response body: ", err.Error())
	}

	var gendata types.GenderData = types.GenderData{}

	err = json.Unmarshal(data, &gendata)
	if err != nil {
		fmt.Println("Error unmarshalling age data: ", err.Error())
	}

	enData.Gender = gendata.Gender
}

func NationEnrichment(url string, wg *sync.WaitGroup, enData *types.MsgEnriched) {
	defer wg.Done()
	httpClient := http.Client{}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		fmt.Println("Error creating request: ", err.Error())
	}

	resp, err := httpClient.Do(req)
	if err != nil {
		fmt.Println("Error doing request: ", err.Error())
	}

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading response body: ", err.Error())
	}

	var nationData types.NationData = types.NationData{}

	err = json.Unmarshal(data, &nationData)
	if err != nil {
		fmt.Println("Error unmarshalling age data: ", err.Error())
	}

	enData.Nationalites = nationData.Nationalities
}
