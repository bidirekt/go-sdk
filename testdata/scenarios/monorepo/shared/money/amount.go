package money

type Amount struct {
	Cents    int64  `json:"cents" bidirekt:"field"`
	Currency string `json:"currency" bidirekt:"field description='ISO 4217 code'"`
}

func (this Amount) Add(other Amount) Amount {
	return Amount{Cents: this.Cents + other.Cents, Currency: this.Currency}
}
