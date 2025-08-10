package domain

type Mandate struct {
	From  string `json:"from"`
	To    string `json:"to"`
	Value int8   `json:"value"`
}

func NewMandate(from, to string, value int8) *Mandate {
	m := new(Mandate)
	m.From = from
	m.To = to
	m.Value = value
	return m
}
