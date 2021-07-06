package dlock

import (
	"fmt"
	"strings"
)

// Validate
type Validate struct {
	err []string
}

// NewValidate new validate
func NewValidate() *Validate {
	return &Validate{
		err: []string{},
	}
}

// StringIsNull is string field is null?
func (c *Validate) StringIsNull(field, fieldName string) *Validate {
	if len(field) <= 0 {
		c.err = append(c.err, fmt.Sprintf("%s is null", fieldName))
	}
	return c
}

// ObjectNull is string field is null?
func (c *Validate) ObjectNull(field interface{}, fieldName string) *Validate {
	if field == nil {
		c.err = append(c.err, fmt.Sprintf("%s is null", fieldName))
	}
	return c
}

// Int64IsNull is int64 field is null
func (c *Validate) Int64IsNull(field int64, fieldName string) *Validate {
	if field <= 0 {
		c.err = append(c.err, fmt.Sprintf("%s is null", fieldName))
	}
	return c
}

// SliceEmpty is slice is null
func (c *Validate) SliceEmpty(field []string, fieldName string) *Validate {
	if len(field) <= 0 {
		c.err = append(c.err, fmt.Sprintf("%s is null", fieldName))
	}
	return c
}

// ToError  all validate error string  to error
func (c *Validate) ToError() error {
	if len(c.err) > 0 {
		return fmt.Errorf(strings.Join(c.err, ";"))
	}
	return nil
}
