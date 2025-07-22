package payments

import (
	"database/sql/driver"
	"fmt"
	"strings"

	"github.com/shopspring/decimal"
)

type Decimal struct {
	decimal.Decimal
}

func (d Decimal) MarshalJSON() ([]byte, error) {
	return []byte(d.Decimal.String()), nil
}

func (d *Decimal) UnmarshalJSON(data []byte) error {
	s := strings.Trim(string(data), `"`)
	dec, err := decimal.NewFromString(s)
	if err != nil {
		return fmt.Errorf("invalid decimal %q: %w", s, err)
	}
	d.Decimal = dec
	return nil
}

func (d *Decimal) Scan(src any) error {
	return d.Decimal.Scan(src)
}

func (d Decimal) Value() (driver.Value, error) {
	return d.Decimal.Value()
}
