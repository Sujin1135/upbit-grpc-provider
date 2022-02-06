package trader

import (
	"testing"
	"upbit-grpc-provider/internal/config"
)

var trader = NewTrader()

func init() {
	config.InitLogger()
	config.InitConfigForTest()
}

func TestTrader_GetAccounts(t *testing.T) {
	accounts := trader.GetAccounts()

	if len(accounts) < 0 {
		t.Fatalf("Failed to get all accounts")
	}
}
