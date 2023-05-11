package products

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockRepository struct {
	mock.Mock
}

func NewMockRepository() MockRepository {
	return MockRepository{}
}

func (r *MockRepository) GetAllBySeller(sellerID string) ([]Product, error) {
	args := r.Mock.Called(sellerID)
	pl := args.Get(0).([]Product)
	err := args.Error(1)

	return pl, err
}

func TestNewService(t *testing.T) {
	//Arrange
	repo := NewMockRepository()
	expected := &service{
		repo: &repo,
	}

	//Act
	sv := NewService(&repo)

	//Assert
	assert.Equal(t, expected, sv)
}

func TestGetAllBySellerService(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		//Arrange
		repo := NewMockRepository()
		sellerId := "1"
		expected := []Product{
			{
				ID:          "1",
				SellerID:    "1",
				Description: "The description",
				Price:       100,
			},
			{
				ID:          "2",
				SellerID:    "1",
				Description: "The description",
				Price:       200,
			},
		}

		repo.On("GetAllBySeller", sellerId).Return([]Product{
			{
				ID:          "1",
				SellerID:    "1",
				Description: "The description",
				Price:       100,
			},
			{
				ID:          "2",
				SellerID:    "1",
				Description: "The description",
				Price:       200,
			},
		}, nil)
		sv := NewService(&repo)

		//Act
		pl, err := sv.GetAllBySeller(sellerId)

		//Assert
		assert.NoError(t, err)
		assert.Equal(t, expected, pl)
		repo.AssertExpectations(t)
	})

	t.Run("Error", func(t *testing.T) {
		//Arrange
		repo := NewMockRepository()
		sellerId := "-1"
		expectedError := errors.New("The ID was not found")

		repo.On("GetAllBySeller", sellerId).Return([]Product{}, errors.New("The ID was not found"))
		sv := NewService(&repo)

		//Act
		pl, err := sv.GetAllBySeller(sellerId)

		//Assert
		assert.EqualError(t, err, expectedError.Error())
		assert.Nil(t, pl)
		repo.AssertExpectations(t)
	})
}
