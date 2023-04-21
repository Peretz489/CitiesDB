//Модуль описаний структур

package structs

type CityInfo struct {
	Id         int    `json:"id"`
	Name       string `json:"name"`
	Region     string `json:"region"`
	District   string `json:"district"`
	Population int    `json:"population"`
	Foundation int    `json:"foundation"`
}

type NewPopulation struct {
	Value int `json:"value"`
}

type StringQuery struct {
	Request string `json:"request"`
}

type Values struct {
	MinValue int `json:"min_value,omitempty"`
	MaxValue int `json:"max_value,omitempty"`
}
