package feishu

import (
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestEmailHelper_Lookup(t *testing.T) {
	logrus.SetLevel(logrus.DebugLevel)
	helper, err := NewEmailHelper(getAppConf())
	require.Nil(t, err)
	openIDs, err := helper.Lookup([]string{"john.xu@cardinfolink.com", "tommy.shang@cardinfolink.com"})
	require.Nil(t, err)
	require.Equal(t, 2, len(openIDs))

	openIDs, err = helper.Lookup([]string{"john.xu@cardinfolink.com", "tommy.shang@cardinfolink.com"})
	require.Nil(t, err)
	require.Equal(t, 2, len(openIDs))
}
