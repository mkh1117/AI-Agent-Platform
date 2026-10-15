package rabbitmq

import (
	"fmt"
	"net"
	"net/url"
	"strings"

	"github.com/mkh1117/agentops/internal/config"
	amqp "github.com/rabbitmq/amqp091-go"
)

type Client struct {
	conn *amqp.Connection

	channel *amqp.Channel
}

func New(cfg config.RabbitMQConfig) (*Client, error) {
	vhost := cfg.VHost
	if vhost == "" {
		vhost = "/"
	}

	path := "/" + strings.TrimPrefix(vhost, "/")

	u := url.URL{
		Scheme: "amqp",
		User:   url.UserPassword(cfg.User, cfg.Password),
		Host:   net.JoinHostPort(cfg.Host, fmt.Sprint(cfg.Port)),
		Path:   path,
	}

	conn, err := amqp.Dial(u.String())
	if err != nil {
		return nil, fmt.Errorf("connect to RabbitMQ: %w", err)
	}

	channel, err := conn.Channel()
	if err != nil {
		_ = conn.Close()
		return nil, fmt.Errorf("open RabbitMQ channel: %w", err)
	}

	return &Client{
		conn:    conn,
		channel: channel,
	}, nil
}

func (c *Client) Channel() *amqp.Channel {
	return c.channel
}

func (c *Client) Close() error {
	var firstErr error

	if c.channel != nil {
		if err := c.channel.Close(); err != nil {
			firstErr = fmt.Errorf("close RabbitMQ channel: %w", err)
		}
	}

	if c.conn != nil {
		if err := c.conn.Close(); err != nil && firstErr == nil {
			firstErr = fmt.Errorf("close RabbitMQ connection: %w", err)
		}
	}

	return firstErr
}