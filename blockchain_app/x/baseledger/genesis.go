package baseledger

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/unibrightio/baseledger/x/baseledger/keeper"
	"github.com/unibrightio/baseledger/x/baseledger/types"
)

// InitGenesis initializes the capability module's state from a provided genesis
// state.
func InitGenesis(ctx sdk.Context, k keeper.Keeper, genState types.GenesisState) {
	// this line is used by starport scaffolding # genesis/module/init
	// Set all the BaseledgerTransaction
	for _, elem := range genState.BaseledgerTransactionList {
		k.SetBaseledgerTransaction(ctx, *elem)
	}

	// Set BaseledgerTransaction count
	k.SetBaseledgerTransactionCount(ctx, genState.BaseledgerTransactionCount)

	// this line is used by starport scaffolding # ibc/genesis/init
}

// ExportGenesis returns the capability module's exported genesis.
func ExportGenesis(ctx sdk.Context, k keeper.Keeper) *types.GenesisState {
	genesis := types.DefaultGenesis()

	// this line is used by starport scaffolding # genesis/module/export
	// Get all BaseledgerTransaction
	BaseledgerTransactionList := k.GetAllBaseledgerTransaction(ctx)
	for _, elem := range BaseledgerTransactionList {
		elem := elem
		genesis.BaseledgerTransactionList = append(genesis.BaseledgerTransactionList, &elem)
	}

	// Set the current count
	genesis.BaseledgerTransactionCount = k.GetBaseledgerTransactionCount(ctx)

	// this line is used by starport scaffolding # ibc/genesis/export

	return genesis
}
