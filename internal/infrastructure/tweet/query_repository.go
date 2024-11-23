package tweet

import (
	"context"
	"database/sql"

	"github.com/jmoiron/sqlx"
	"github.com/nhutHao02/social-network-common-service/utils/logger"
	"github.com/nhutHao02/social-network-tweet-service/internal/domain/interface/tweet"
	"github.com/nhutHao02/social-network-tweet-service/internal/domain/model"
	"go.uber.org/zap"
)

type tweetQueryRepository struct {
	db *sqlx.DB
}

// GetLoveTweetsByUserID implements tweet.TweetQueryRepository.
func (repo *tweetQueryRepository) GetLoveTweetsByUserID(ctx context.Context, req model.GetLoveTweetsByUserIDReq) ([]model.GetTweetsRes, uint64, error) {
	var res []model.GetTweetsRes
	query := `select t.ID ,
					t.Content ,
					t.UserID
				from lovetweet lt 
				left join tweet t
				on lt.TweetID = t.ID 
				where lt.UserID = :UserID and lt.DeletedAt is null and t.DeletedAt is null
				order by t.CreatedAt desc 
				limit :Limit Offset :Offset`
	params := map[string]interface{}{
		"UserID": req.UserID,
		"Limit":  req.Limit,
		"Offset": (req.Page - 1) * req.Limit,
	}

	queryString, args, err := repo.db.BindNamed(query, params)
	if err != nil {
		logger.Error("tweetQueryRepository-GetLoveTweetsByUserID: bindName for query error", zap.Error(err))
		return res, 0, err
	}

	err = repo.db.SelectContext(ctx, &res, queryString, args...)
	if err != nil && err != sql.ErrNoRows {
		logger.Error("tweetQueryRepository-GetLoveTweetsByUserID: get tweets error", zap.Error(err))
		return res, 0, err
	}

	// count total tweets
	var count uint64
	queryCount := `select count(*) from lovetweet lt 
				left join tweet t
				on lt.TweetID = t.ID 
				where lt.UserID = ? and lt.DeletedAt is null and t.DeletedAt is null`
	err = repo.db.GetContext(ctx, &count, queryCount, req.UserID)
	if err != nil {
		logger.Error("tweetQueryRepository-GetLoveTweetsByUserID: count total tweet error", zap.Error(err))
		return res, 0, nil
	}

	// get tweet statistics and user action to tweet
	for index, item := range res {
		// pass statistics
		statistics, _ := GetTweetStatistics(ctx, repo.db, item.ID)
		res[index].Statistics = &statistics

		// pass actions
		action, _ := GetUserActionWithTweet(ctx, repo.db, item.ID, int(req.UserID))
		res[index].UserAction = &action
	}

	return res, count, nil
}

// GetTweets implements tweet.TweetQueryRepository.
func (repo *tweetQueryRepository) GetTweets(ctx context.Context, req model.GetTweetsReq) ([]model.GetTweetsRes, uint64, error) {
	var res []model.GetTweetsRes
	query := `select t.ID ,
					t.Content ,
					t.UserID
				from tweet t 
				where t.DeletedAt is null 
				order by t.CreatedAt desc 
				limit :Limit Offset :Offset`
	params := map[string]interface{}{
		"Limit":  req.Limit,
		"Offset": (req.Page - 1) * req.Limit,
	}

	queryString, args, err := repo.db.BindNamed(query, params)
	if err != nil {
		logger.Error("tweetQueryRepository-GetTweets: bindName for query error", zap.Error(err))
		return res, 0, err
	}
	err = repo.db.SelectContext(ctx, &res, queryString, args...)
	if err != nil && err != sql.ErrNoRows {
		logger.Error("tweetQueryRepository-GetTweets: get tweets error", zap.Error(err))
		return res, 0, err
	}

	// count total tweets
	queryCount := `select count(*) from tweet tw where tw.DeletedAt is null`
	var count uint64
	err = repo.db.GetContext(ctx, &count, queryCount)
	if err != nil {
		logger.Error("tweetQueryRepository-GetTweets: count total tweet error", zap.Error(err))
		return res, 0, nil
	}

	// get tweet statistics and user action to tweet
	for index, item := range res {
		// pass statistics
		statistics, _ := GetTweetStatistics(ctx, repo.db, item.ID)
		res[index].Statistics = &statistics

		// pass actions
		action, _ := GetUserActionWithTweet(ctx, repo.db, item.ID, int(req.UserID))
		res[index].UserAction = &action
	}

	return res, count, nil

}

func GetUserActionWithTweet(ctx context.Context, db *sqlx.DB, tweetID int, userID int) (model.UserAction, error) {
	var res model.UserAction
	query := `SELECT
				(SELECT IF(COUNT(l.ID) > 0, 1, 0) FROM lovetweet l WHERE l.TweetID = :TweetID and l.UserID = :UserID) AS Love,
				(SELECT IF(COUNT(b.ID) > 0, 1, 0) FROM bookmarktweet b WHERE b.TweetID = :TweetID and b.UserID = :UserID) AS Bookmark,
				(SELECT IF(COUNT(r.ID) > 0, 1, 0) FROM reposttweet r WHERE r.TweetID = :TweetID and r.UserID = :UserID) AS Repost;`

	params := map[string]interface{}{
		"TweetID": tweetID,
		"UserID":  userID,
	}

	queryString, args, err := db.BindNamed(query, params)
	if err != nil {
		logger.Warn("tweetQueryRepository-GetUserActionWithTweet: bind name query error", zap.Error(err))
		return res, err
	}

	err = db.GetContext(ctx, &res, queryString, args...)
	if err != nil {
		logger.Warn("tweetQueryRepository-GetUserActionWithTweet: get user action with tweet error", zap.Error(err))
		return res, err
	}
	return res, nil
}

func GetTweetStatistics(ctx context.Context, db *sqlx.DB, tweetID int) (model.Statistics, error) {
	var res model.Statistics
	query := `SELECT
				(SELECT COUNT(distinct l.UserID) FROM lovetweet l WHERE l.TweetID = ?) AS TotalLove,
				(SELECT COUNT(distinct tc.UserID) FROM tweetcomment tc WHERE tc.TweetID = ?) AS TotalComment,
				(SELECT COUNT(distinct b.UserID) FROM bookmarktweet b WHERE b.TweetID = ?) AS TotalBookmark,
				(SELECT COUNT(distinct r.UserID) FROM reposttweet r WHERE r.TweetID = ?) AS TotalRepost;`
	err := db.GetContext(ctx, &res, query, tweetID, tweetID, tweetID, tweetID)
	if err != nil {
		logger.Warn("tweetQueryRepository-GetTweetStatistics: get tweet statistics error", zap.Error(err))
		return res, err
	}
	return res, nil
}

// GetTweetByUserID implements tweet.TweetQueryRepository.
func (repo *tweetQueryRepository) GetTweetByUserID(ctx context.Context, req model.GetTweetByUserReq) ([]model.GetTweetByUserRes, uint64, error) {
	var res []model.GetTweetByUserRes
	query := `select t.ID ,
					t.Content ,
					t.UserID
				from tweet t 
				where t.UserID = :UserID and t.DeletedAt is null 
				order by t.CreatedAt desc 
				limit :Limit Offset :Offset`
	params := map[string]interface{}{
		"UserID": req.UserID,
		"Limit":  req.Limit,
		"Offset": (req.Page - 1) * req.Limit,
	}

	queryString, args, err := repo.db.BindNamed(query, params)
	if err != nil {
		logger.Error("repo-GetTweetByUserID: bindName for query error", zap.Error(err))
		return res, 0, err
	}
	err = repo.db.SelectContext(ctx, &res, queryString, args...)
	if err != nil && err != sql.ErrNoRows {
		logger.Error("repo-GetTweetByUserID: get tweets error", zap.Error(err))
		return res, 0, err
	}

	queryCount := `select count(*) from tweet tw where tw.UserID = ? and tw.DeletedAt is null`
	var count uint64
	err = repo.db.GetContext(ctx, &count, queryCount, req.UserID)
	if err != nil {
		logger.Error("repo-GetTweetByUserID: count total tweet error", zap.Error(err))
		return res, 0, nil
	}

	// get tweet statistics and user action to tweet
	for index, item := range res {
		// pass statistics
		statistics, _ := GetTweetStatistics(ctx, repo.db, item.ID)
		res[index].Statistics = &statistics

		// pass actions
		action, _ := GetUserActionWithTweet(ctx, repo.db, item.ID, req.UserID)
		res[index].UserAction = &action
	}

	return res, count, nil
}

func NewTweetQueryRepository(db *sqlx.DB) tweet.TweetQueryRepository {
	return &tweetQueryRepository{db: db}
}
