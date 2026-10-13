package ai

import (
	"context"
)

type ModelProvider interface{
	Generate(cnx context.Context,message Message)(string,error)
}


