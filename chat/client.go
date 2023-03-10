package chat

import (
	"fmt"
	"github.com/ejedavy/scalable-chat-server/utils"
	"github.com/go-redis/redis/v7"
	"github.com/gorilla/websocket"
)

type Client struct {
	conn          *websocket.Conn
	redisClient   *redis.Client
	username      string
	pubSub        *redis.PubSub
	listening     bool
	stopListening chan struct{}
	messageChan   chan redis.Message
	config        utils.Config
}

func NewClient(username string, conn *websocket.Conn, redisClient *redis.Client, config utils.Config) (*Client, error) {
	client := Client{
		conn:          conn,
		redisClient:   redisClient,
		username:      username,
		pubSub:        nil,
		listening:     false,
		stopListening: make(chan struct{}),
		messageChan:   make(chan redis.Message),
		config:        config,
	}
	err := client.connect()
	if err != nil {
		return nil, err
	}
	return &client, nil
}

func (c *Client) connect() error {
	channels := make([]string, 0, 10)
	initialChannels, err := c.redisClient.SMembers(c.config.GeneralChannels).Result()
	if err != nil {
		return err
	}
	channels = append(channels, initialChannels...)
	userChannelsKey := fmt.Sprintf(c.config.UserChannelsFmt, c.username)
	subscriptions, err := c.redisClient.SMembers(userChannelsKey).Result()
	if err != nil {
		return err
	}
	channels = append(channels, subscriptions...)
	if len(channels) == 0 {
		return fmt.Errorf("no channel to subscribe user %s\n", c.username)
	}

	if c.pubSub != nil {
		err = c.pubSub.Unsubscribe()
		if err != nil {
			return err
		}
		err = c.pubSub.Close()
		if err != nil {
			return err
		}
	}
	if c.listening {
		c.stopListening <- struct{}{}
	}

	c.pubSub = c.redisClient.Subscribe(channels...)
	go func() {
		fmt.Println("Listening for", c.username, "has started")
		c.listening = true
	loop:
		for {
			select {
			case msg := <-c.pubSub.Channel():
				c.messageChan <- *msg
			case <-c.stopListening:
				fmt.Printf("Listening for %s has been stopped\n", c.username)
				break loop
			}
		}
	}()

	return nil
}

func (c *Client) Subscribe(channel string) error {
	userChannelsKey := fmt.Sprintf(c.config.UserChannelsFmt, c.username)
	if c.redisClient.SIsMember(userChannelsKey, channel).Val() {
		return nil
	}
	c.redisClient.SAdd(userChannelsKey, channel)
	return c.connect()
}

func (c *Client) Unsubscribe(channel string) error {
	userChannelsKey := fmt.Sprintf(c.config.UserChannelsFmt, c.username)
	if !c.redisClient.SIsMember(userChannelsKey, channel).Val() {
		return nil
	}
	c.redisClient.SRem(userChannelsKey, channel)
	return c.connect()
}

func (c *Client) Disconnect() error {
	if c.pubSub != nil {
		err := c.pubSub.Unsubscribe()
		if err != nil {
			return err
		}
		err = c.pubSub.Close()
		if err != nil {
			return err
		}
	}
	if c.listening {
		c.stopListening <- struct{}{}
	}
	close(c.messageChan)
	return nil
}
