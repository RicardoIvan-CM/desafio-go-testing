package products

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetAllBySellerRepository(t *testing.T) {
	//Arrange
	repo := NewRepository()
	sellerId := "1"
	expected := []Product{{
		ID:          "mock",
		SellerID:    "FEX112AC",
		Description: "generic product",
		Price:       123.55,
	}}

	//Act
	pl, err := repo.GetAllBySeller(sellerId)

	//Assert
	assert.NoError(t, err)
	assert.Equal(t, expected, pl)
}
