package uuid

import (
	"database/sql/driver"

	"github.com/google/uuid"
)

func (u UUID) Value() (driver.Value, error) {
	impl := uuid.UUID(u)
	return impl.Value(), nil
}

func (u *UUID) Scan(src interface{}) error {
	var impl uuid.UUID
	err := impl.Scan(src)
	if err != nil {
		return err
	}
	*u = UUID(impl)
	return nil
}
