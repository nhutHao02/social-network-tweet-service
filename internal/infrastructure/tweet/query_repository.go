package tweet

import (
	"github.com/jmoiron/sqlx"
	"github.com/nhutHao02/social-network-tweet-service/internal/domain/interface/tweet"
)

type tweetQueryRepository struct {
	db *sqlx.DB
}

func NewTweetQueryRepository(db *sqlx.DB) tweet.TweetQueryRepository {
	return &tweetQueryRepository{db: db}
}
