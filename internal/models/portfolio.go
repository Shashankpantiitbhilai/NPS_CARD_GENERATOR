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

// Validate ensures all portfolio data is valid
func (p *UserPortfolio) Validate() error {
	// Could add validation logic here, such as:
	// - Ensuring allocations sum to 100%
	// - Validating color formats
	// - Checking for empty scheme names
	return nil
}