// Package notary defines all relevant functionality for a Notary actor
// within a sharded Ethereum blockchain.
package notary

import (
	"fmt"

	"github.com/ethereum/go-ethereum/ethdb"
	"github.com/ethereum/go-ethereum/log"
	"github.com/ethereum/go-ethereum/sharding"
	"github.com/ethereum/go-ethereum/sharding/mainchain"
	"github.com/ethereum/go-ethereum/sharding/params"
)

// Notary holds functionality required to run a collation notary
// in a sharded system. Must satisfy the Service interface defined in
// sharding/service.go.
type Notary struct {
	config       *params.ShardConfig
	smcClient    *mainchain.SMCClient
	shardp2p     sharding.ShardP2P
	shardChainDb ethdb.Database
}

// NewNotary creates a new notary instance.
func NewNotary(config *params.ShardConfig, smcClient *mainchain.SMCClient, shardp2p sharding.ShardP2P, shardChainDb ethdb.Database) (*Notary, error) {
	return &Notary{config, smcClient, shardp2p, shardChainDb}, nil
}

// Start the main routine for a notary.
func (n *Notary) Start() {
	log.Info("Starting notary service")
	go n.notarizeCollations()
}

// Stop the main loop for notarizing collations.
func (n *Notary) Stop() error {
	log.Info("Stopping notary service")
	return nil
}

func (n *Notary) notarizeCollations() {
	// TODO: handle this better through goroutines. Right now, these methods
	// are blocking.
	if n.smcClient.DepositFlag() {
		if err := joinNotaryPool(n.config, n.smcClient); err != nil {
			log.Error(fmt.Sprintf("Could not fetch current block number: %v", err))
			return
		}
	}

	if err := subscribeBlockHeaders(n.smcClient); err != nil {
		log.Error(fmt.Sprintf("Could not fetch current block number: %v", err))
		return
	}
}