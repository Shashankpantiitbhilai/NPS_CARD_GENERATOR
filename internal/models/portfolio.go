package models

// UserPortfolio represents the user's NPS portfolio data.
type UserPortfolio struct {
	Username   string   `json:"username"`
	XIRR       float64  `json:"xirr"`
	XIRRPeriod int      `json:"xirrPeriod"`
	Schemes    []Scheme `json:"schemes"`
}

// Scheme represents an individual scheme in the portfolio.
type Scheme struct {
	Name       string `json:"name"`
	Allocation int    `json:"allocation"` // percentage
	Color      string `json:"color"`      // hex code or named color
}


func (p *UserPortfolio) Validate() error {
	
	return nil
}