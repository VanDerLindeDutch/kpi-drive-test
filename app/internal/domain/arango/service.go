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
	"strconv"
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

func (s *Service) Process() error {
	ctx, _ := context.WithTimeout(context.Background(), time.Minute*2)
	data, err := s.getData(ctx)
	if err != nil {
		return err
	}
	for _, v := range data {
		err := s.saveData(ctx, v)
		if err != nil {
			log.Println(err)
		}
	}
	return nil

}

func (s *Service) saveData(ctx context.Context, row eventResponseBodyRow) error {
	form := url.Values{}
	form.Add("period_start", row.Params.Period.Start)
	form.Add("period_end", row.Params.Period.End)
	form.Add("period_key", row.Params.Period.TypeKey)
	form.Add("indicator_to_mo_id", strconv.Itoa(row.Params.IndicatorToMoId))
	form.Add("indicator_to_mod_fact_id", strconv.Itoa(row.Params.IndicatorToMoId))
	form.Add("fact_time", "2023-08-11")
	tag := superTagBody{
		Tag: tagBody{
			Id:           row.Author.UserId,
			Name:         row.Author.UserName,
			Key:          "client",
			ValuesSource: 0,
		},
		Value: row.Author.UserName,
	}
	marshaledTag, err := json.Marshal(tag)
	if err != nil {
		return err
	}
	form.Add("supertags", string(marshaledTag))

	cookie, err := s.authorize(ctx)
	if err != nil {
		return err
	}
	resp := map[string]interface{}{}
	_, err = s.httpClient.PostFormDataCookie(ctx, "facts/save_fact", nil, form, &resp, cookie)
	for s2, i := range resp {
		fmt.Println(s2, i)
	}

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
	fmt.Println(string(jsonBody))
	if err != nil {
		return nil, err
	}
	fmt.Println(newEventRequestBody(field, []string{"time"}, "DESC", 10))
	resp := eventResponseBody{}
	_, err = s.httpClient.GetJsonCookie(ctx, "events", nil, jsonBody, &resp, cookie)
	for s2, s3 := range resp.Data.Rows {
		fmt.Println(s2, s3)
	}
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
	cookies, i, err := s.httpClient.CookieAuthorize(context.Background(), "auth/login", nil, form, nil)
	if err != nil {
		return nil, err
	}
	if i != 200 || len(cookies) < 1 {
		return nil, fmt.Errorf("unauthorized")
	}
	s.cache.Set("cookies", cookies[0], cache.DefaultExpiration)
	return cookies[0], nil

}
