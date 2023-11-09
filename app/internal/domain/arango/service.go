package arango

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/patrickmn/go-cache"
	"kpi-drive-test/app/internal/config"
	"kpi-drive-test/app/internal/pkg/http_client"
	"log"
	"net/http"
	"net/url"
	"sync"
	"time"
)

type Service struct {
	cfg        *config.Config
	httpClient *http_client.Client
	cache      *cache.Cache
}

func NewService(cfg *config.Config) *Service {
	return &Service{cfg: cfg, httpClient: http_client.NewClient("https://development.kpi-drive.ru/_api/", nil), cache: cache.New(10*time.Minute, 10*time.Minute)}
}

func (s *Service) FetchDataAndSave() error {
	ctx, _ := context.WithTimeout(context.Background(), time.Minute*2)
	data, err := s.getData(ctx)
	if err != nil {
		return err
	}
	wg := sync.WaitGroup{}
	wg.Add(len(data))
	for _, v := range data {
		go func(v eventResponseBodyRow) {
			defer wg.Done()

			err := s.saveData(ctx, v)
			if err != nil {
				log.Println("ERROR :", err)
			}

		}(v)

	}
	wg.Wait()
	return nil

}

func (s *Service) saveData(ctx context.Context, row eventResponseBodyRow) error {

	factTime := arangoDate{time: &row.Time}
	//fmt.Println(today.String())

	tag := superTagBody{
		Tag: tagBody{
			Id:           row.Author.UserId,
			Name:         "Клиент",
			Key:          "client",
			ValuesSource: 0,
		},
		Value: row.Author.UserName,
	}
	marshaledTag, err := json.Marshal(tag)
	if err != nil {
		return err
	}

	params, _ := json.Marshal(row.Params)

	form := url.Values{}
	form.Add("period_start", row.Params.Period.Start.String())
	form.Add("period_end", row.Params.Period.End.String())
	form.Add("period_key", row.Params.Period.TypeKey)
	form.Add("indicator_to_mo_id", "315914")
	form.Add("indicator_to_mod_fact_id", "0")
	form.Add("fact_time", factTime.String())
	form.Add("is_plan", "0")
	form.Add("value", "1")
	form.Add("supertags", string(marshaledTag))
	form.Add("auth_user_id", "40")
	form.Add("comment", string(params))

	cookie, err := s.authorize(ctx)
	if err != nil {
		return err
	}
	resp := map[string]interface{}{}
	_, err = s.httpClient.PostFormDataCookie(ctx, "facts/save_fact", nil, form, &resp, cookie)
	respData := resp["DATA"]
	parsedData := respData.(map[string]interface{})

	marshal, _ := json.MarshalIndent(row, "", "    ")
	log.Println(fmt.Sprintf("EVENT RESPONSE:\n %s", string(marshal)))
	log.Println(fmt.Sprintf("Indicator id: %v", parsedData["indicator_to_mo_fact_id"]))

	return nil

}

func (s *Service) getData(ctx context.Context) ([]eventResponseBodyRow, error) {

	cookie, err := s.authorize(ctx)
	if err != nil {
		return nil, err
	}
	field := eventRequestBodyFilterField{
		Key:    "type",
		Sign:   "LIKE",
		Values: []string{"MATRIX_REQUEST"},
	}
	reqBody := newEventRequestBody(field, []string{"time"}, "DESC", 10)

	jsonBody, err := json.Marshal(reqBody)
	if err != nil {
		return nil, err
	}
	resp := eventResponseBody{}
	_, err = s.httpClient.GetJsonCookie(ctx, "events", nil, jsonBody, &resp, cookie)

	return resp.Data.Rows, nil
}

func (s *Service) authorize(ctx context.Context) (*http.Cookie, error) {
	cachedCookies, ok := s.cache.Get("cookies")
	if ok {
		return cachedCookies.(*http.Cookie), nil
	}
	form := url.Values{}
	form.Add("login", s.cfg.Kpi.Username)
	form.Add("password", s.cfg.Kpi.Password)
	cookies, i, err := s.httpClient.CookieAuthorize(ctx, "auth/login", nil, form, nil)
	if err != nil {
		return nil, err
	}
	if i != 200 || len(cookies) < 1 {
		return nil, fmt.Errorf("unauthorized")
	}
	s.cache.Set("cookies", cookies[0], cache.DefaultExpiration)
	return cookies[0], nil

}
