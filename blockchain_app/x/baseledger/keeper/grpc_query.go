package keeper

import (
	"github.com/unibrightio/baseledger/x/baseledger/types"
)

var _ types.QueryServer = Keeper{}
