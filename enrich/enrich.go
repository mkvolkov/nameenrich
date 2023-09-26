package enrich

import (
	"context"
	"io"
	"nameenrich/types"
	"net/http"
	"sync"
	"time"

	"github.com/goccy/go-json"
)

func Enrichment(baseMsg *types.MsgBase, enrMsg *types.MsgEnriched) (err error) {
	enrMsg.Name = baseMsg.Name
	enrMsg.Surname = baseMsg.Surname
	enrMsg.Patronymic = baseMsg.Patronymic

	urlAge := "https://api.agify.io/?name=" + baseMsg.Name
	urlGender := "https://api.genderize.io/?name=" + baseMsg.Name
	urlNation := "https://api.nationalize.io/?name=" + baseMsg.Name

	var wg sync.WaitGroup

	var ageCh chan error = make(chan error)
	var genderCh chan error = make(chan error)
	var nationCh chan error = make(chan error)

	wg.Add(1)
	go AgeEnrichment(urlAge, &wg, enrMsg, ageCh)

	wg.Add(1)
	go GenderEnrichment(urlGender, &wg, enrMsg, genderCh)

	wg.Add(1)
	go NationEnrichment(urlNation, &wg, enrMsg, nationCh)

	wg.Wait()

	select {
	case err1 := <-ageCh:
		err = err1
	case err2 := <-genderCh:
		err = err2
	case err3 := <-nationCh:
		err = err3
	default:
		err = nil
	}

	return err
}

func AgeEnrichment(url string, wg *sync.WaitGroup, enData *types.MsgEnriched, errCh chan error) {
	defer wg.Done()

	httpClient := http.Client{}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		errCh <- err
		return
	}

	resp, err := httpClient.Do(req)
	if err != nil {
		errCh <- err
		return
	}

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		errCh <- err
		return
	}

	var agedata types.AgeData = types.AgeData{}

	err = json.Unmarshal(data, &agedata)
	if err != nil {
		errCh <- err
		return
	}

	enData.Age = agedata.Age
}

func GenderEnrichment(url string, wg *sync.WaitGroup, enData *types.MsgEnriched, errCh chan error) {
	defer wg.Done()
	httpClient := http.Client{}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		errCh <- err
		return
	}

	resp, err := httpClient.Do(req)
	if err != nil {
		errCh <- err
		return
	}

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		errCh <- err
		return
	}

	var gendata types.GenderData = types.GenderData{}

	err = json.Unmarshal(data, &gendata)
	if err != nil {
		errCh <- err
		return
	}

	enData.Gender = gendata.Gender
}

func NationEnrichment(url string, wg *sync.WaitGroup, enData *types.MsgEnriched, errCh chan error) {
	defer wg.Done()
	httpClient := http.Client{}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		errCh <- err
		return
	}

	resp, err := httpClient.Do(req)
	if err != nil {
		errCh <- err
		return
	}

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		errCh <- err
		return
	}

	var nationData types.NationData = types.NationData{}

	err = json.Unmarshal(data, &nationData)
	if err != nil {
		errCh <- err
		return
	}

	enData.Nationalites = nationData.Nationalities
}
