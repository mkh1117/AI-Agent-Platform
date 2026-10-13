package ai

import (
    "testing"
	"encoding/json"

    "github.com/stretchr/testify/require"
)

func TestJsonUnmarshal(t *testing.T){
	data:=[]byte(`{"role":"user","content":"سلام"}`)
	var message Message
 
	err:=json.Unmarshal(data,&message)

	require.NoError(t,err)
	
	require.Equal(t,"سلام",message.Content)

}
func TestRoleValidate(t *testing.T){
	var Y Message
	var X Message
	var Z Message

	Y.Content="help"
	Y.Role=Role("banana")

	X.Content="hello"
	X.Role=RoleUser
	
	Z.Content="   "
	Z.Role=RoleAssistant

	ErrorY:=Y.Validate()

	ErrorX:=X.Validate()

	ErrorZ:=Z.Validate()

	require.Error(t,ErrorY)
	require.NoError(t,ErrorX)
	require.Error(t,ErrorZ)
}