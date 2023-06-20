package tokentype

import (
	"context"
	"math/big"

	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/sirupsen/logrus"
)

type Client struct {
	ctx    context.Context
	logger *logrus.Entry

	ethClient *ethclient.Client
	chainId   *big.Int

	apiKey  string
	isDebug bool
}

func NewClient(httpRpc string) (*Client, error) {
	var err error

	client := &Client{
		ctx:    context.Background(),
		logger: logrus.WithField("tokentype", "apitypes"),
	}
	client.ethClient, err = ethclient.DialContext(client.ctx, httpRpc)
	if err != nil {
		client.logger.WithError(err).Error("DialContext ethclient failed")
		return nil, err
	}

	client.chainId, err = client.ethClient.ChainID(client.ctx)
	if err != nil {
		client.logger.WithError(err).Error("ethClient get_chainId failed")
		return nil, err
	}
	return client, nil
}

func (c *Client) WithAPIKey(apiKey string) {
	c.apiKey = apiKey
}

func (c *Client) WithDebug(isDebug bool) {
	c.isDebug = isDebug
}
