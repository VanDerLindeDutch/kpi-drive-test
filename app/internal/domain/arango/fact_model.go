package arango

type superTagBody struct {
	Tag   tagBody `json:"tag"`
	Value string  `json:"value"`
}

type tagBody struct {
	Id           int    `json:"id"`
	Name         string `json:"name"`
	Key          string `json:"key"`
	ValuesSource int    `json:"values_source"`
}
