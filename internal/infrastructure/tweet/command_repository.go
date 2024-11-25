package tweet

import (
	"context"

	"github.com/jmoiron/sqlx"
	"github.com/nhutHao02/social-network-common-service/utils/logger"
	"github.com/nhutHao02/social-network-tweet-service/internal/domain/interface/tweet"
	"github.com/nhutHao02/social-network-tweet-service/internal/domain/model"
	"github.com/nhutHao02/social-network-tweet-service/pkg/constants"
	"go.uber.org/zap"
)

type tweetCommandRepository struct {
	db        *sqlx.DB
	queryRepo tweet.TweetQueryRepository
}

// PostComment implements tweet.TweetCommandRepository.
func (repo *tweetCommandRepository) PostComment(ctx context.Context, params map[string]interface{}) (model.OutgoingMessageWSRes, error) {
	query := `INSERT INTO sntweetservice.tweetcomment
			(Description, UserID, TweetID)
			VALUES(:Description, :UserID, :TweetID);`
	_, err := repo.db.NamedExecContext(ctx, query, params)
	if err != nil {
		logger.Error("tweetCommandRepository-PostComment: Error when Execute query ", zap.Error(err))
	}
	outGoingMessageWSRes, err := repo.queryRepo.GetNewCommentTweetByUserIDAndTweetID(ctx, params)
	if err == nil {
		return outGoingMessageWSRes, err
	}
	return outGoingMessageWSRes, nil
}

func getQueryDeleteActionTweets(req model.ActionTweetReq) string {
	queryUpdate := ``
	switch req.Action {
	case constants.Love:
		queryUpdate = `UPDATE lovetweet `
	case constants.Bookmark:
		queryUpdate = `UPDATE bookmarktweet `
	default:
		queryUpdate = `UPDATE reposttweet `
	}

	querySet := `SET DeletedAt=current_timestamp() `
	queryClauses := `WHERE UserID = :UserID and TweetID = :TweetID`

	return queryUpdate + querySet + queryClauses
}

// DeleteActionTweetsByUserID implements tweet.TweetCommandRepository.
func (repo *tweetCommandRepository) DeleteActionTweetsByUserID(ctx context.Context, req model.ActionTweetReq) (bool, error) {
	// check exist tweet
	_, err := repo.queryRepo.ExistedTweet(ctx, int64(req.TweetID))
	if err != nil {
		logger.Error("tweetCommandRepository-DeleteActionTweetsByUserID: error when check Exist tweet", zap.Error(err))
		return false, err
	}
	query := getQueryDeleteActionTweets(req)

	_, err = repo.db.NamedExecContext(ctx, query, req)
	if err != nil {
		logger.Error("tweetCommandRepository-DeleteActionTweetsByUserID: error when Execute context", zap.Error(err))
		return false, nil
	}
	return true, nil
}

func getQueryActionTweets(req model.ActionTweetReq) string {
	queryInsert := ``
	switch req.Action {
	case constants.Love:
		queryInsert = `INSERT INTO lovetweet `
	case constants.Bookmark:
		queryInsert = `INSERT INTO bookmarktweet `
	default:
		queryInsert = `INSERT INTO reposttweet `
	}

	queryFields := `(UserID, TweetID)
					VALUES(:UserID, :TweetID);`

	return queryInsert + queryFields
}

// ActionTweetsByUserID implements tweet.TweetCommandRepository.
func (repo *tweetCommandRepository) ActionTweetsByUserID(ctx context.Context, req model.ActionTweetReq) (bool, error) {
	// check exist tweet
	_, err := repo.queryRepo.ExistedTweet(ctx, int64(req.TweetID))
	if err != nil {
		logger.Error("tweetCommandRepository-ActionTweetsByUserID: error when check Exist tweet", zap.Error(err))
		return false, err
	}
	query := getQueryActionTweets(req)

	_, err = repo.db.NamedExecContext(ctx, query, req)
	if err != nil {
		logger.Error("tweetCommandRepository-ActionTweetsByUserID: error when Execute context", zap.Error(err))
		return false, nil
	}
	return true, nil

}

func savePostVideo(ctx context.Context, tx *sqlx.Tx, params map[string]interface{}) error {
	query := `INSERT INTO sntweetservice.tweetimage
				(Url, TweetID)
				VALUES(:UrlVideo, :TweetID);`

	_, err := tx.NamedExecContext(ctx, query, params)
	if err != nil {
		logger.Error("tweetCommandRepository-savePostVideo: error when save Post Video", zap.Error(err))
		return err
	}
	return nil
}

func savePostImg(ctx context.Context, tx *sqlx.Tx, params map[string]interface{}) error {
	query := `INSERT INTO sntweetservice.tweetvideo
				(Url, TweetID)
				VALUES(:UrlImg, :TweetID);`

	_, err := tx.NamedExecContext(ctx, query, params)
	if err != nil {
		logger.Error("tweetCommandRepository-savePostImg: error when save Post Image", zap.Error(err))
		return err
	}
	return nil
}

// PostTweet implements tweet.TweetCommandRepository.
func (repo *tweetCommandRepository) PostTweet(ctx context.Context, req model.PostTweetReq) (bool, error) {
	tx, err := repo.db.BeginTxx(ctx, nil)
	if err != nil {
		logger.Error("tweetCommandRepository-PostTweet: fail to begin transaction ", zap.Error(err))
	}
	defer func() {
		if err != nil {
			tx.Rollback()
			logger.Error("tweetCommandRepository-PostTweet: transaction rollback due to error ", zap.Error(err))
			return
		}
		err = tx.Commit()
		if err != nil {
			logger.Error("tweetCommandRepository-PostTweet: failed to commit transaction ", zap.Error(err))
		}

	}()

	query := `INSERT INTO sntweetservice.tweet
				(Content, UserID)
				VALUES(:Content, :UserID);`

	result, err := tx.NamedExecContext(ctx, query, req)
	if err != nil {
		return false, err
	}
	tweetID, err := result.LastInsertId()
	if err != nil {
		logger.Error("tweetCommandRepository-PostTweet error when get tweetID", zap.Error(err))
		return false, err
	}
	if req.UrlImg != nil {
		err = savePostImg(ctx, tx, map[string]interface{}{
			"TweetID": tweetID,
			"UrlImg":  req.UrlImg,
		})

		if err != nil {
			return false, err
		}
	}
	if req.UrlVideo != nil {
		err = savePostVideo(ctx, tx, map[string]interface{}{
			"TweetID":  tweetID,
			"UrlVideo": req.UrlVideo,
		})

		if err != nil {
			return false, err
		}
	}

	return true, nil
}

func NewTweetCommandRepository(db *sqlx.DB, queryRepo tweet.TweetQueryRepository) tweet.TweetCommandRepository {
	return &tweetCommandRepository{db: db, queryRepo: queryRepo}
}
