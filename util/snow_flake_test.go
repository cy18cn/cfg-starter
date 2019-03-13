package util

import "testing"

func TestSnowFlake_NextId(t *testing.T) {
	sf, err := NewSnowFlake(&SFSettings{})
	if err != nil {
		t.Error(err)
	}

	var id uint64
	id, err = sf.NextId()
	if err != nil {
		t.Error(err)
	}
	t.Log(id)
}
