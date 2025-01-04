package tweet

import (
	"context"
	"database/sql"

	"github.com/jmoiron/sqlx"
	resError "github.com/nhutHao02/social-network-common-service/utils/error"
	"github.com/nhutHao02/social-network-common-service/utils/logger"
	"github.com/nhutHao02/social-network-tweet-service/internal/domain/interface/tweet"
	"github.com/nhutHao02/social-network-tweet-service/internal/domain/model"
	"github.com/nhutHao02/social-network-tweet-service/pkg/constants"
	"go.uber.org/zap"
)

type tweetQueryRepository struct {
	db *sqlx.DB
}

// GetAuthorIDOfTweet implements tweet.TweetQueryRepository.
func (repo *tweetQueryRepository) GetAuthorIDOfTweet(ctx context.Context, tweetId int64) (int64, error) {
	var userID int64
	query := `select t.UserID 
				from tweet t 
				where t.ID = ? and t.DeletedAt is null `

	if err := repo.db.GetContext(ctx, &userID, query, tweetId); err != nil {
		logger.Error("tweetQueryRepository-GetAuthorIDOfTweet: get userID error", zap.Error(err))
		return userID, err
	}
	return userID, nil
}

// GetTweetComments implements tweet.TweetQueryRepository.
func (repo *tweetQueryRepository) GetTweetComments(ctx context.Context, req model.TweetCommentReq) ([]model.TweetCommentRes, uint64, error) {
	var res []model.TweetCommentRes
	query := `select t.ID ,
						t.Description ,
						t.UserID ,
						t.CreatedAt 
				from tweetcomment t
				where t.TweetID = :TweetID and t.DeletedAt is null 
				order by t.CreatedAt desc 
				limit :Limit Offset :Offset`
	params := map[string]interface{}{
		"TweetID": req.TweetID,
		"Limit":   req.Limit,
		"Offset":  (req.Page - 1) * req.Limit,
	}
	queryString, args, err := repo.db.BindNamed(query, params)
	if err != nil {
		logger.Error("tweetQueryRepository-GetTweetComments: bindName for query error", zap.Error(err))
		return res, 0, err
	}

	err = repo.db.SelectContext(ctx, &res, queryString, args...)
	if err != nil && err != sql.ErrNoRows {
		logger.Error("tweetQueryRepository-GetTweetComments: get tweets error", zap.Error(err))
		return res, 0, err
	}

	queryCount := `select count(*) from tweetcomment tc 
					where tc.TweetID = ? and tc.DeletedAt is null `
	var count uint64
	err = repo.db.GetContext(ctx, &count, queryCount, req.TweetID)
	if err != nil {
		logger.Error("tweetQueryRepository-GetTweetComments: count total comment error", zap.Error(err))
		return res, 0, nil
	}

	return res, count, nil
}

// GetNewCommentTweetByUserIDAndTweetID implements tweet.TweetQueryRepository.
func (repo *tweetQueryRepository) GetNewCommentTweetByUserIDAndTweetID(ctx context.Context, params map[string]interface{}) (model.OutgoingMessageWSRes, error) {
	var res model.OutgoingMessageWSRes
	query := `select tc.ID as 'ID' ,
					tc.UserID as 'UserID',
					tc.Description as 'Content' ,
					tc.CreatedAt as 'Timestamp' 
				from tweetcomment tc 
				where tc.UserID = :UserID and tc.TweetID = :TweetID and tc.DeletedAt is null 
				order by tc.ID desc 
				limit 1`
	queryString, args, err := repo.db.BindNamed(query, params)
	if err != nil {
		logger.Error("tweetQueryRepository-GetNewCommentTweetByUserIDAndTweetID: Error when BindName query", zap.Error(err))
		return res, err
	}
	if err = repo.db.GetContext(ctx, &res, queryString, args...); err != nil {
		logger.Error("tweetQueryRepository-GetNewCommentTweetByUserIDAndTweetID: Error when GetContext", zap.Error(err))
		return res, err
	}
	return res, nil
}

// ExistedTweet implements tweet.TweetQueryRepository.
func (repo *tweetQueryRepository) ExistedTweet(ctx context.Context, tweetId int64) (bool, error) {
	var count int64
	query := `select count(*) from tweet t where t.ID = ? and t.DeletedAt is null`
	err := repo.db.GetContext(ctx, &count, query, tweetId)
	if err != nil {
		logger.Error("ExistedTweet: error ", zap.Error(err))
		return false, err
	}
	if count < 1 {
		return false, resError.NewResError(nil, "Tweet not exist")
	}
	return true, nil
}

func getQueryActionTweetsByUserID(actionTweet constants.ActionTweet) string {
	query := `select t.ID ,
					t.Content ,
					t.UserID,
					t.CreatedAt `
	queryJoin := ``
	switch actionTweet {
	case constants.Love:
		queryJoin = `from lovetweet t1 
					left join tweet t 
					on t1.TweetID = t.ID `
	case constants.Bookmark:
		queryJoin = `from bookmarktweet t1 
					left join tweet t 
					on t1.TweetID = t.ID `
	case constants.Post:
		queryJoin = `from tweet t `
	default:
		queryJoin = `from reposttweet t1 
					left join tweet t 
					on t1.TweetID = t.ID `
	}
	queryClauses := ``
	if actionTweet == constants.Post {
		queryClauses = `where t.UserID = :UserID and t.DeletedAt is null
		order by t.CreatedAt desc 
		limit :Limit Offset :Offset`
	} else {
		queryClauses = `where t1.UserID = :UserID and t1.DeletedAt is null and t.DeletedAt is null
		order by t.CreatedAt desc 
		limit :Limit Offset :Offset`
	}
	return query + queryJoin + queryClauses
}
func getQueryCountActionTweet(actionTweet constants.ActionTweet) string {
	query := `select count(*) `
	queryJoin := ``
	switch actionTweet {
	case constants.Love:
		queryJoin = `from lovetweet t1 
					left join tweet t 
					on t1.TweetID = t.ID `
	case constants.Bookmark:
		queryJoin = `from bookmarktweet t1 
					left join tweet t 
					on t1.TweetID = t.ID `
	default:
		queryJoin = `from reposttweet t1 
					left join tweet t 
					on t1.TweetID = t.ID `
	}
	queryClauses := `where t1.UserID = ? and t1.DeletedAt is null and t.DeletedAt is null`
	return query + queryJoin + queryClauses
}

// GetActionTweetsByUserID implements tweet.TweetQueryRepository.
func (repo *tweetQueryRepository) GetActionTweetsByUserID(ctx context.Context, req model.GetActionTweetsByUserIDReq) ([]model.GetTweetsRes, uint64, error) {
	var res []model.GetTweetsRes
	query := getQueryActionTweetsByUserID(req.Action)
	params := map[string]interface{}{
		"UserID": req.UserID,
		"Limit":  req.Limit,
		"Offset": (req.Page - 1) * req.Limit,
	}

	queryString, args, err := repo.db.BindNamed(query, params)
	if err != nil {
		logger.Error("tweetQueryRepository-GetActionTweetsByUserID: bindName for query error", zap.Error(err))
		return res, 0, err
	}

	err = repo.db.SelectContext(ctx, &res, queryString, args...)
	if err != nil && err != sql.ErrNoRows {
		logger.Error("tweetQueryRepository-GetActionTweetsByUserID: get tweets error", zap.Error(err))
		return res, 0, err
	}

	// count total tweets
	var count uint64
	queryCount := getQueryCountActionTweet(req.Action)
	err = repo.db.GetContext(ctx, &count, queryCount, req.UserID)
	if err != nil {
		logger.Error("tweetQueryRepository-GetActionTweetsByUserID: count total tweet error", zap.Error(err))
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
					img.Url as 'UrlImg',
					video.Url as 'UrlVideo',
					t.UserID,
					t.CreatedAt 
				from tweet t 
				left join tweetimage img on t.ID = img.TweetID
				left join tweetvideo video on t.ID  = video.TweetID 
				where t.DeletedAt is null and img.DeletedAt is null and video.DeletedAt is null 
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
				(SELECT IF(COUNT(l.ID) > 0, 1, 0) FROM lovetweet l WHERE l.TweetID = :TweetID and l.UserID = :UserID and l.DeletedAt is null) AS Love,
				(SELECT IF(COUNT(b.ID) > 0, 1, 0) FROM bookmarktweet b WHERE b.TweetID = :TweetID and b.UserID = :UserID and b.DeletedAt is null) AS Bookmark,
				(SELECT IF(COUNT(r.ID) > 0, 1, 0) FROM reposttweet r WHERE r.TweetID = :TweetID and r.UserID = :UserID and r.DeletedAt is null) AS Repost;`

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
				(SELECT COUNT(distinct l.UserID) FROM lovetweet l WHERE l.TweetID = ? and l.DeletedAt is null) AS TotalLove,
				(SELECT COUNT(distinct tc.UserID) FROM tweetcomment tc WHERE tc.TweetID = ? and tc.DeletedAt is null) AS TotalComment,
				(SELECT COUNT(distinct b.UserID) FROM bookmarktweet b WHERE b.TweetID = ? and b.DeletedAt is null) AS TotalBookmark,
				(SELECT COUNT(distinct r.UserID) FROM reposttweet r WHERE r.TweetID = ? and r.DeletedAt is null) AS TotalRepost;`
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
