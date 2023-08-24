package serenity

import (
	"time"

	"github.com/sagernet/serenity/libsubscription"
	"github.com/sagernet/sing-box/experimental/libbox"
	"github.com/sagernet/sing/common"
	E "github.com/sagernet/sing/common/exceptions"
	"github.com/sagernet/sing/common/logger"
)

type SubscriptionClient struct {
	logger        logger.Logger
	ticker        *time.Ticker
	subscriptions []*SubscriptionOptions
	close         chan struct{}
}

func NewSubscriptionClient(logger logger.Logger, subscriptions []*SubscriptionOptions) *SubscriptionClient {
	var interval time.Duration
	for _, subscription := range subscriptions {
		if interval == 0 || time.Duration(subscription.UpdateInterval) < interval {
			interval = time.Duration(subscription.UpdateInterval)
		}
	}
	if interval == 0 {
		interval = 5 * time.Minute
	}
	return &SubscriptionClient{
		logger:        logger,
		ticker:        time.NewTicker(interval),
		subscriptions: subscriptions,
		close:         make(chan struct{}),
	}
}

func (c *SubscriptionClient) Subscriptions() []*SubscriptionOptions {
	return common.Filter(c.subscriptions, func(it *SubscriptionOptions) bool {
		return it.LastUpdate != time.Time{}
	})
}

func (c *SubscriptionClient) Start() error {
	err := c.Update()
	if err != nil {
		return err
	}
	go c.loopUpdate()
	return nil
}

func (c *SubscriptionClient) Close() error {
	c.ticker.Stop()
	select {
	case <-c.close:
	default:
		close(c.close)
	}
	return nil
}

func (c *SubscriptionClient) loopUpdate() {
	for {
		select {
		case <-c.ticker.C:
		case <-c.close:
			return
		}
		err := c.Update()
		if err != nil {
			c.logger.Error(err)
		}
	}
}

func (c *SubscriptionClient) Update() error {
	httpClient := libbox.NewHTTPClient()
	httpClient.ModernTLS()
	defer httpClient.Close()
	for i, subscription := range c.subscriptions {
		err := c.update(httpClient, subscription)
		if err != nil {
			c.logger.Error(E.Cause(err, "update subscription[", i, "]: ", subscription.URL))
		}
		c.logger.Info("updated subscription[", i, "]: ", len(subscription.ServerCache), " servers")
	}
	return nil
}

func (c *SubscriptionClient) update(httpClient libbox.HTTPClient, subscription *SubscriptionOptions) error {
	request := httpClient.NewRequest()
	err := request.SetURL(subscription.URL)
	if err != nil {
		return err
	}
	if subscription.UserAgent != "" {
		request.SetUserAgent(subscription.UserAgent)
	} else {
		request.SetUserAgent("ClashForAndroid/serenity")
	}
	response, err := request.Execute()
	if err != nil {
		return err
	}
	content, err := response.GetContentString()
	if err != nil {
		return err
	}
	servers, err := libsubscription.ParseSubscription(content)
	if err != nil {
		return err
	}
	subscription.LastUpdate = time.Now()
	subscription.ServerCache = servers
	return nil
}
