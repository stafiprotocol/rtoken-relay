package solana

import (
	"context"
	"fmt"
	"time"

	"github.com/ChainSafe/log15"
	"github.com/stafiprotocol/rtoken-relay/chains"
	"github.com/stafiprotocol/rtoken-relay/core"
	solClient "github.com/stafiprotocol/solana-go-sdk/client"
)

var (
	BlockRetryInterval = time.Second * 3
	BlockRetryLimit    = 400
	EraInterval        = time.Minute * 1
)

//listen event or block update from solana
type listener struct {
	name   string
	symbol core.RSymbol
	conn   *Connection
	router chains.Router
	log    log15.Logger
	stop   <-chan int
	sysErr chan<- error
}

func NewListener(name string, symbol core.RSymbol, conn *Connection, log log15.Logger, stop <-chan int, sysErr chan<- error) *listener {
	return &listener{
		name:   name,
		symbol: symbol,
		conn:   conn,
		log:    log,
		stop:   stop,
		sysErr: sysErr,
	}
}

func (l *listener) setRouter(r chains.Router) {
	l.router = r
}

func (l *listener) start() error {

	go func() {
		err := l.pollBlocks()
		if err != nil {
			l.log.Error("Polling blocks failed", "err", err)
			panic(err)
		}
	}()

	return nil
}

func (l *listener) pollBlocks() error {
	var retry = BlockRetryLimit
	ticker := time.NewTicker(EraInterval)
	l.updateEra()
	defer ticker.Stop()
	for {
		select {
		case <-l.stop:
			return ErrorTerminated
		case <-ticker.C:
			if retry <= 0 {
				return fmt.Errorf("poolBlocks reach retry limit ,symbol: %s", l.symbol)
			}

			err := l.updateEra()
			if err != nil {
				retry--
				continue
			}
			retry = BlockRetryLimit
		}
	}
}

func (l *listener) updateEra() error {
	client := l.conn.GetQueryClient()

	epochInfo, err := client.GetEpochInfo(context.Background(), solClient.CommitmentFinalized)
	if err != nil {
		return err
	}
	currentEra := uint32(epochInfo.Epoch)

	msg := &core.Message{Destination: core.RFIS, Reason: core.NewEra, Content: currentEra}
	l.submitMessage(msg, nil)
	return nil
}

// submitMessage inserts the chainId into the msg and sends it to the router
func (l *listener) submitMessage(m *core.Message, err error) {
	if err != nil {
		l.log.Error("Critical error before sending message", "err", err)
		return
	}
	m.Source = l.symbol
	err = l.router.Send(m)
	if err != nil {
		l.log.Error("failed to send message", "err", err, "msg", m)
	}
}
