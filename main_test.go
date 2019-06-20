package main

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSilly(t *testing.T) {
	asrt, rq := assert.New(t), require.New(t)

	var s struct{ S *SillyTime }
	asrt.NoError(json.Unmarshal([]byte(`{"s": 1}`), &s))
	rq.NotNil(s.S)

	rq.Equal(time.Unix(1, 0), s.S.Time)
}
