package config

import (
	"github.com/davecgh/go-spew/spew"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestLoadFromFile(t *testing.T) {
	conf, err := Load("../config.example.yml")
	require.Nil(t, err)
	spew.Dump(conf)
}
