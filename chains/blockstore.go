package chains

import (
	"errors"
	"github.com/stafiprotocol/chainbridge/utils/blockstore"
	"github.com/stafiprotocol/chainbridge/utils/msg"
	"github.com/stafiprotocol/rtoken-relay/core"
	"strconv"
)

func StartBlock(cfg *core.ChainConfig, latestBlock uint64, bs *blockstore.Blockstore, relayer string) (uint64, error) {
	if cfg.LatestBlockFlag {
		return latestBlock, nil
	}

	blk := parseStartBlock(cfg)

	bsCfg := cfg.Opts["blockstorePath"]
	bsPath, ok := bsCfg.(string)
	if !ok {
		return 0, errors.New("blockstorePath not string")
	}

	var err error
	bs, err = blockstore.NewBlockstore(bsPath, msg.ChainId(100), relayer)
	if err != nil {
		return 0, err
	}

	blk, err = checkBlockstore(bs, blk)
	if err != nil {
		return 0, err
	}

	return blk, nil
}

func parseStartBlock(cfg *core.ChainConfig) uint64 {
	if blk, ok := cfg.Opts["startBlock"]; ok {
		blkStr, ok := blk.(string)
		if !ok {
			panic("block not string")
		}
		res, err := strconv.ParseUint(blkStr, 10, 32)
		if err != nil {
			panic(err)
		}
		return res
	}
	return 0
}

// checkBlockstore queries the blockstore for the latest known block. If the latest block is
// greater than startBlock, then the latest block is returned, otherwise startBlock is.
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
