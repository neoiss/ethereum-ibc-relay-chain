package ethereum

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/hyperledger-labs/yui-relayer/core"
)

var _ core.ChainConfig = (*ChainConfig)(nil)

func (c ChainConfig) Build() (core.Chain, error) {
	return NewChain(c)
}

func (c ChainConfig) IBCAddress() common.Address {
	return common.HexToAddress(c.IbcAddress)
}

func (c ChainConfig) ICS20BankAddress() common.Address {
	return common.HexToAddress(c.Ics20BankAddress)
}

func (c ChainConfig) ICS20TransferBankAddress() common.Address {
	return common.HexToAddress(c.Ics20TransferBankAddress)
}
