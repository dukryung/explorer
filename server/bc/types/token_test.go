package types

import (
	"encoding/json"
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/hessegg/nikto/x/bankz/types"
	"github.com/stretchr/testify/assert"
)

func TestJSON(t *testing.T) {
	pt := types.PlatformToken()
	zero := sdk.NewInt64Coin(pt.GetDenom(), int64(0))
	token := Token{
		Token: pt,
		Coin:  &zero,
	}
	bz, err := json.Marshal(token)
	assert.NoError(t, err)
	t.Log(string(bz))
	t.Fail()
}
