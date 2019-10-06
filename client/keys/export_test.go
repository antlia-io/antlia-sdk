package keys

import (
	"testing"

	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"

	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/tests"
)

func Test_runExportCmd(t *testing.T) {
	runningOnServer := isRunningOnServer()
	exportKeyCommand := exportKeyCommand()
	mockIn, _, _ := tests.ApplyMockIO(exportKeyCommand)

	// Now add a temporary keybase
	kbHome, cleanUp := tests.NewTestCaseDir(t)
	defer cleanUp()
	viper.Set(flags.FlagHome, kbHome)

	// create a key
	kb := NewKeyring(mockIn)
	if !runningOnServer {
		defer func() {
			kb.Delete("keyname1", "", false)
		}()
	}

	if runningOnServer {
		mockIn.Reset("testpass1\ntestpass1\n")
	}
	_, err := kb.CreateAccount("keyname1", tests.TestMnemonic, "", "123456789", 0, 0)
	assert.NoError(t, err)

	// Now enter password
	if runningOnServer {
		mockIn.Reset("123456789\ntestpass1\n")
	} else {
		mockIn.Reset("123456789\n123456789\n")
	}
	assert.NoError(t, runExportCmd(exportKeyCommand, []string{"keyname1"}))
}
