package tweet

import (
	"github.com/jmoiron/sqlx"
	"github.com/nhutHao02/social-network-tweet-service/internal/domain/interface/tweet"
)

type tweetCommandRepository struct {
	db *sqlx.DB
}

func NewTweetCommandRepository(db *sqlx.DB) tweet.TweetCommandRepository {
	return &tweetCommandRepository{db: db}
}
