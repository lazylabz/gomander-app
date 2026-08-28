package test

import (
	"context"

	"github.com/stretchr/testify/mock"
)

type MockAppControl struct {
	mock.Mock
}

func (m *MockAppControl) CloseApp(ctx context.Context) {
	m.Called(ctx)
}
