package chains

import (
	"errors"
	"github.com/stafiprotocol/chainbridge/utils/blockstore"
	"github.com/stafiprotocol/chainbridge/utils/msg"
	"strconv"
)

func NewBlockstore(bsCfg interface{}, relayer string) (*blockstore.Blockstore, error) {
	bsPath, ok := bsCfg.(string)
	if !ok {
		return nil, errors.New("blockstorePath not string")
	}

	//todo change chainId for different rToken
	return blockstore.NewBlockstore(bsPath, msg.ChainId(100), relayer)
}

func StartBlock(bs *blockstore.Blockstore, blkCfg interface{}) (uint64, error) {
	blk := parseStartBlock(blkCfg)
	return checkBlockstore(bs, blk)
}

func parseStartBlock(blkCfg interface{}) uint64 {
	blkStr, ok := blkCfg.(string)
	if !ok {
		panic("blkCfg not string")
	}

	res, err := strconv.ParseUint(blkStr, 10, 32)
	if err != nil {
		panic(err)
	}
	return res
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
