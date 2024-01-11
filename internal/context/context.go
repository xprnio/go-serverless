package context

import (
	"github.com/xprnio/go-serverless/internal/functions"
	"github.com/google/uuid"
)

type Context struct {
	Function    *functions.Function
	RequestId   string
	RequestBody []byte
}

func New(request []byte) (*Context, error) {
	requestId, err := uuid.NewRandom()
	if err != nil {
		return nil, err
	}

	return &Context{
		RequestId:   requestId.String(),
		RequestBody: request,
	}, nil
}
