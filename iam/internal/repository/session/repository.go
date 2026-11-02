package session

import "github.com/redis/go-redis/v9"

type repository struct {
	client redis.Client
}

func New(client redis.Client) *repository {
	return &repository{
		client: client,
	}
}
