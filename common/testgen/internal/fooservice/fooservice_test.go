package fooservice

import (
	mock_fooservice "github.com/fidesy/sdk/common/testgen/internal/fooservice/mocks"
	"github.com/golang/mock/gomock"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type FooServiceShould struct {
	suite.Suite
}

func TestFooServiceShould(t *testing.T) {
	suite.Run(t, new(FooServiceShould))
}

func (s *FooServiceShould) TestCall_AllValidParams_NoError() {
	ctrl := gomock.NewController(s.T())

	externalServiceMock := mock_fooservice.NewMockExternalService(ctrl)
	externalServiceMock.EXPECT().
		Call().
		Return(nil)

	service := New(
		WithExternalService(externalServiceMock),
	)

	err := service.Call()

	require.NoError(s.T(), err)
}
