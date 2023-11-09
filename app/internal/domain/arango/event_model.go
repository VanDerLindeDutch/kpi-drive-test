package arango

import (
	"strings"
	"time"
)

type eventRequestBody struct {
	Filter struct {
		Field eventRequestBodyFilterField `json:"field"`
	} `json:"filter"`
	Sort struct {
		Fields    []string `json:"fields"`
		Direction string   `json:"direction"`
	} `json:"sort"`

	Limit int `json:"limit"`
}

type eventRequestBodyFilterField struct {
	Key    string   `json:"key"`
	Sign   string   `json:"sign"`
	Values []string `json:"values"`
}

func newEventRequestBody(field eventRequestBodyFilterField, sortFields []string, sortDirection string, limit int) eventRequestBody {
	reqBody := eventRequestBody{}

	reqBody.Filter.Field = field
	reqBody.Sort.Fields = sortFields
	reqBody.Sort.Direction = sortDirection
	reqBody.Limit = limit
	return reqBody
}

type arangoDate struct {
	time *time.Time
}

func (e *arangoDate) UnmarshalJSON(b []byte) error {
	value := strings.Trim(string(b), `"`) //get rid of "
	if value == "" || value == "null" {
		return nil
	}

	t, err := time.Parse("2006-01-02", value) //parse time
	if err != nil {
		return err
	}
	e.time = &t //set result using the pointer
	return nil
}

func (e arangoDate) MarshalJSON() ([]byte, error) {
	return []byte(`"` + e.time.Format("2006-01-02") + `"`), nil
}

func (e arangoDate) String() string {
	return e.time.Format("2006-01-02")
}

type eventResponseBody struct {
	Data struct {
		Rows []eventResponseBodyRow `json:"rows"`
	} `json:"DATA"`

	Limit int `json:"limit"`
}

type eventResponseBodyRow struct {
	Id     string                  `json:"_id"`
	Group  string                  `json:"group"`
	Type   string                  `json:"type"`
	Time   time.Time               `json:"time"`
	Author eventResponseBodyAuthor `json:"author"`
	Params struct {
		IndicatorToMoId int                     `json:"indicator_to_mo_id"`
		Period          eventResponseBodyPeriod `json:"period"`
	} `json:"params"`
}

type eventResponseBodyAuthor struct {
	UserId   int    `json:"user_id"`
	UserName string `json:"user_name"`
	MoId     int    `json:"mo_id"`
}

type eventResponseBodyPeriod struct {
	End     arangoDate `json:"end"`
	Start   arangoDate `json:"start"`
	TypeId  int        `json:"type_id"`
	TypeKey string     `json:"type_key"`
}
