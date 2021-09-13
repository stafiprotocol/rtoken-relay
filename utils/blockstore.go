package utils

import (
	"github.com/stafiprotocol/chainbridge/utils/blockstore"
	"github.com/stafiprotocol/chainbridge/utils/crypto/sr25519"
	"github.com/stafiprotocol/chainbridge/utils/keystore"
	"github.com/stafiprotocol/chainbridge/utils/msg"
	"github.com/stafiprotocol/go-substrate-rpc-client/signature"
	"strconv"
)

func BlockstoreAndKeyring(from, keystorePath, blockstorePath, startBlock string) (*signature.KeyringPair, *blockstore.Blockstore, uint64, error) {
	kp, err := keystore.KeypairFromAddress(from, keystore.SubChain, keystorePath, false)
	if err != nil {
		return nil, nil, 0, err
	}

	krp := kp.(*sr25519.Keypair).AsKeyringPair()

	// Attempt to load latest block
	bs, err := blockstore.NewBlockstore(blockstorePath, msg.ChainId(1), kp.Address())
	if err != nil {
		return nil, nil, 0, err
	}
	blk := parseStartBlock(startBlock)
	blk, err = checkBlockstore(bs, blk)
	if err != nil {
		return nil, nil, 0, err
	}

	return krp, bs, blk, nil
}

func checkBlockstore(bs *blockstore.Blockstore, startBlock uint64) (uint64, error) {
	latestBlock, err := bs.TryLoadLatestBlock()
	if err != nil {
		return 0, err
	}

	if latestBlock.Uint64() > startBlock {
		return latestBlock.Uint64(), nil
	} else {
		return startBlock, nil
	}
}

func parseStartBlock(startBlock string) uint64 {
	//if blk, ok := cfg.Opts[""]; ok {
	blk, err := strconv.ParseUint(startBlock, 10, 32)
	if err != nil {
		panic(err)
	}
	return blk
}
