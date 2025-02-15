import { useNavigate } from "react-router";
import { TTweet } from "../../types/tweet";
import { UserIcon } from "../user-icon/user-icon";
import { LikeButton } from "../like-button/like-button";
import { ReplyButton } from "../reply-button/reply-button";
import { RetweetButton } from "../retweet-button/retweet-button";
import "./tweet-list-item.css";
import { displayTime } from "../../utils/display-time";
import RetweetItem from "../retweet-item";

interface TweetListItemProps {
  tweet: TTweet;
  refetch: () => void;
}

export const TweetListItem = ({ tweet, refetch }: TweetListItemProps) => {
  const navigate = useNavigate();

  return (
    <div
      key={tweet.id}
      className="TweetListItem"
      onClick={() => navigate(`/${tweet.user.display_id}/status/${tweet.id}`)}
    >
      <div className="TweetListItem__user__icon__wrapper">
        <UserIcon user={tweet.user} size={40} />
      </div>
      <div className="TweetListItem__content__wrapper">
        <div className="TweetListItem__content__user__wrapper">
          <span
            className="TweetListItem__content__user__name"
            onClick={(e) => {
              e.stopPropagation();
              navigate(`/${tweet.user.display_id}`);
            }}
          >
            {tweet.user.name}
          </span>
          <span className="TweetListItem__content__user__displayID">
            @{tweet.user.display_id}
          </span>
          <span className="TweetListItem__content__user__createdAt">
            ・{displayTime(tweet.created_at)}
          </span>
        </div>
        <div className="TweetListItem__content__content">{tweet.content}</div>
        {tweet.retweet && <RetweetItem retweet={tweet.retweet} />}
        <div className="TweetListItem__content__interactions">
          <ReplyButton
            tweet={tweet}
            className="TweetListItem__content__interactions__reply"
          />
          <RetweetButton
            tweet={tweet}
            className="TweetListItem__content__interactions__retweet"
          />
          <LikeButton
            tweet={tweet}
            refetch={refetch}
            className="TweetListItem__content__interactions__like"
          />
        </div>
      </div>
    </div>
  );
};
