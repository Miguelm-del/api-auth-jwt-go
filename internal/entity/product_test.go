package entity

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewProduct(t *testing.T) {
	p, err := NewProduct("Tea", 10.0)
	assert.Nil(t, err)
	assert.NotNil(t, p)
	assert.NotEmpty(t, p.ID)
	assert.Equal(t, "Tea", p.Name)
	assert.Equal(t, 10, p.Price)
}

func TestProduct_WhenNameIsRequired(t *testing.T) {
	p, err := NewProduct("", 10.0)
	assert.Nil(t, p)
	assert.Equal(t, ErrNameIsRequired, err)
}

func TestProduct_WhenPriceIsRequired(t *testing.T) {
	p, err := NewProduct("Beans", 0.0)
	assert.Nil(t, p)
	assert.Equal(t, ErrPriceIsRequired, err)
}

func TestProduct_WhenPriceIsInvalid(t *testing.T) {
	p, err := NewProduct("Beans", -10.0)
	assert.Nil(t, p)
	assert.Equal(t, ErrInvalidPrice, err)
}

func TestProduct_Validate(t *testing.T) {
	p, err := NewProduct("Coffee", 10.0)
	assert.Nil(t, err)
	assert.NotNil(t, p)
	assert.Nil(t, p.Validate())
}
