package common

import (
	"github.com/rs/xid"
	"testing"
)

type TestApi struct {
	testService *TestService `inject:"testService" single:"true"`
	stringVar   string
}

type TestService struct {
	intVar   int
	floatVar float32
}

func TestNewField(t *testing.T) {
	a := make([]string, 5)
	casted(a, t)
}

func casted(i interface{}, t *testing.T) {
	switch i.(type) {
	case []string:
		t.Logf("[]string: %v", i)
	}
}

func TestNewDefinition(t *testing.T) {
	id:= xid.New()

	t.Log(id.String())
	t.Log(id.String())
}
