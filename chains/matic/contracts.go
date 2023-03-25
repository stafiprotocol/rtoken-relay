package matic

import (
	"errors"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/stafiprotocol/rtoken-relay/bindings/MaticToken"
	stake_portal "github.com/stafiprotocol/rtoken-relay/bindings/StakeERC20Portal"
	stake_portal_with_rate "github.com/stafiprotocol/rtoken-relay/bindings/StakeERC20PortalWithRate"
	"github.com/stafiprotocol/rtoken-relay/bindings/StakeManager"
	stake_portal_rate "github.com/stafiprotocol/rtoken-relay/bindings/StakePortalRate"
)

func initStakeManager(stakeManagerCfg interface{}, conn *ethclient.Client) (*StakeManager.StakeManager, common.Address, error) {
	stakeManagerAddr, ok := stakeManagerCfg.(string)
	if !ok {
		return nil, ZeroAddress, errors.New("StakeManagerContract not ok")
	}
	addr := common.HexToAddress(stakeManagerAddr)
	manager, err := StakeManager.NewStakeManager(common.HexToAddress(stakeManagerAddr), conn)
	if err != nil {
		return nil, ZeroAddress, err
	}

	return manager, addr, nil
}

func initStakePortal(stakeManagerCfg interface{}, conn *ethclient.Client) (*stake_portal.StakeERC20Portal, error) {
	stakePortalAddr, ok := stakeManagerCfg.(string)
	if !ok {
		return nil, errors.New("StakeERC20Portal not ok")
	}
	stakePortal, err := stake_portal.NewStakeERC20Portal(common.HexToAddress(stakePortalAddr), conn)
	if err != nil {
		return nil, err
	}

	return stakePortal, nil
}
func initPolygonStakePortalRate(stakeManagerCfg interface{}, conn *ethclient.Client) (*stake_portal_rate.StakePortalRate, error) {
	stakePortalAddr, ok := stakeManagerCfg.(string)
	if !ok {
		return nil, errors.New("StakeERC20Portal not ok")
	}
	stakePortal, err := stake_portal_rate.NewStakePortalRate(common.HexToAddress(stakePortalAddr), conn)
	if err != nil {
		return nil, err
	}

	return stakePortal, nil
}

func initStakePortalWithRate(stakeManagerCfg interface{}, conn *ethclient.Client) (*stake_portal_with_rate.StakeERC20PortalWithRate, error) {
	stakePortalAddr, ok := stakeManagerCfg.(string)
	if !ok {
		return nil, errors.New("StakeERC20PortalWithRate not ok")
	}
	stakePortal, err := stake_portal_with_rate.NewStakeERC20PortalWithRate(common.HexToAddress(stakePortalAddr), conn)
	if err != nil {
		return nil, err
	}

	return stakePortal, nil
}

func initMaticToken(maticTokenCfg interface{}, conn *ethclient.Client) (*MaticToken.MaticToken, common.Address, error) {
	maticTokenAddr, ok := maticTokenCfg.(string)
	if !ok {
		return nil, ZeroAddress, errors.New("MaticTokenContract not ok")
	}
	addr := common.HexToAddress(maticTokenAddr)
	matic, err := MaticToken.NewMaticToken(addr, conn)
	if err != nil {
		return nil, ZeroAddress, err
	}

	return matic, addr, nil
}

func initMultisend(multiSendCfg interface{}, conn *ethclient.Client) (common.Address, error) {
	multiSendAddr, ok := multiSendCfg.(string)
	if !ok {
		return ZeroAddress, errors.New("MultiSendContract not ok")
	}
	return common.HexToAddress(multiSendAddr), nil
}
