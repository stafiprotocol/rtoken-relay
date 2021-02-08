package service

import (
	"context"
	"fmt"
	"github.com/stafiprotocol/rtoken-relay/conn"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/ChainSafe/log15"
	"github.com/stafiprotocol/rtoken-relay/config"
	"github.com/stafiprotocol/rtoken-relay/substrate"
	"github.com/stafiprotocol/rtoken-relay/utils"
)

type ChainType string

const (
	sub = ChainType("substrate")
)

var (
	validators map[conn.RSymbol]conn.Validator
)

func Start(cfg *config.Config, log log15.Logger) {
	ctx, cancel := context.WithCancel(context.Background())
	sysErr := make(chan error)

	wg := sync.WaitGroup{}
	wg.Add(1)
	go func(ctx context.Context) {
		mc := cfg.MainConf
		sblk, ok := mc.Opts["startBlock"]
		if !ok {
			sblk = "0"
		}
		krp, bs, blk, err := utils.BlockstoreAndKeyring(mc.From, mc.KeystorePath, sblk)
		if err != nil {
			sysErr <- fmt.Errorf("BlockstoreAndKeyring error: %s", err)
			return
		}

		sc, err := substrate.NewSarpcClient(cfg.MainConf.Endpoint, log.New("Sarpc-"+cfg.MainConf.Name))
		if err != nil {
			sysErr <- fmt.Errorf("NewSarpcClient error: %s", err)
			return
		}

		gc, err := substrate.NewGsrpcClient(ctx, cfg.MainConf.Endpoint, krp, log.New("Gsrpc-"+cfg.MainConf.Name))
		if err != nil {
			sysErr <- fmt.Errorf("NewGsrpcClient error: %s", err)
			return
		}

		validators, err = Validators(cfg, log)
		if err != nil {
			sysErr <- err
			return
		}

		listener := NewListener(ctx, sc, gc, bs, blk, validators, sysErr, log)
		err = listener.Start()
		if err != nil {
			sysErr <- fmt.Errorf("listener start error: %s", err)
			return
		}
	}(ctx)

	sigc := make(chan os.Signal, 1)
	signal.Notify(sigc, syscall.SIGINT, syscall.SIGTERM)
	defer signal.Stop(sigc)

	// Block here and wait for a signal
	select {
	case err := <-sysErr:
		log.Error("FATAL ERROR. Shutting down.", "err", err)
	case <-sigc:
		log.Warn("Interrupt received, shutting down now.")
	}

	cancel()
}

func Validators(cfg *config.Config, log log15.Logger) (map[conn.RSymbol]conn.Validator, error) {
	vals := make(map[conn.RSymbol]conn.Validator)
	//
	//for _, chainConf := range cfg.OtherConfs {
	//	ctype := ChainType(chainConf.Type)
	//	switch ctype {
	//	case sub:
	//		sc, err := substrate.NewSarpcClient(chainConf.Endpoint, log.New(chainConf.Name))
	//		if err != nil {
	//			return nil, fmt.Errorf("NewSarpcClient error: %s for chain: %s", err, chainConf.Name)
	//		}
	//		vals[conn.RSymbol(chainConf.Rsymbol)] = sc
	//	}
	//}

	return vals, nil
}
