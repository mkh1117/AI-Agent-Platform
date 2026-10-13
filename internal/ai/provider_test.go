package ai

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

type MockProvider struct{}


func (p MockProvider)Generate(cnx context.Context,message Message)(string ,error){
	return "hello",nil

}

func TestMockProvider(t *testing.T){
	var provider ModelProvider
	  
	provider=MockProvider{}

	result,err:=provider.Generate(context.Background(),Message{})

	require.NoError(t,err)
	require.Equal(t,"hello",result)

}