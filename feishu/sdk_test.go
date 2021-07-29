package feishu

import (
	"github.com/davecgh/go-spew/spew"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestSdk_BatchGetID(t *testing.T) {
	logrus.SetLevel(logrus.DebugLevel)
	app := getAppConf()
	sdk := NewSDK(app.ID, app.Secret)
	resp, err := sdk.TenantAccessToken()
	require.Nil(t, err)
	spew.Dump(resp)

	sdk.token = resp.TenantAccessToken

	openIDs, err := sdk.BatchGetID([]string{"john.xu@cardinfolink.com", "tommy.shang@cardinfolink.com"})
	require.Nil(t, err)
	spew.Dump(openIDs)
	require.Equal(t, 2, len(openIDs))
}
