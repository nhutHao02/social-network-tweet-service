package constants

type ActionTweet string

const (
	Love     ActionTweet = "Love"
	Bookmark ActionTweet = "Bookmark"
	Repost   ActionTweet = "Repost"
	Comment  ActionTweet = "Comment"
)

func (a ActionTweet) IsValid() bool {
	switch a {
	case Love, Bookmark, Repost:
		return true
	default:
		return false
	}
}
