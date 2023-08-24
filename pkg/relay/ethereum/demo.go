package ethereum

import (
	"context"
	"math/big"
	"strings"

	"github.com/ethereum/go-ethereum/common"
)

func (c *Chain) QueryBankBalance(address common.Address, id string) (*big.Int, error) {
	idLower := strings.ToLower(id)
	return c.ics20Bank.BalanceOf(c.CallOpts(context.Background(), -1), address, idLower)
}
