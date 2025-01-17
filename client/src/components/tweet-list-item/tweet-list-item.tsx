import { useNavigate } from "react-router";
import { TTweet } from "../../types/tweet";
import { serverFetch } from "../../utils/fetch";
import { UserIcon } from "../user-icon/user-icon";
import { AiOutlineRetweet } from "react-icons/ai";
import { FaRegComment } from "react-icons/fa6";
import { BsHeart } from "react-icons/bs";
import React from "react";
import "./tweet-list-item.css";
import { displayTime } from "../../utils/display-time";
import RetweetItem from "../retweet-item";

interface TweetListItemProps {
  tweet: TTweet;
  refetch: () => void; // エラーしないようにnullを追加
}

export const TweetListItem = ({ tweet, refetch }: TweetListItemProps) => {
  // refetchはとりあえず残したいけど今は使わない。エラーだけ外したかった
  const navigate = useNavigate();
  const likeTweet = async (e: React.MouseEvent) => {
    e.stopPropagation();
    const res = await serverFetch(`/api/like`, {
      method: "POST",
      body: JSON.stringify({
        tweet_id: tweet.id,
      }),
      headers: {
        "Content-Type": "application/json",
      },
    });
    if (res.ok) {
      refetch();
      console.log("いいねしました");
    }
  };

  return (
    <div
      key={tweet.id}
      className="TweetListItem"
      onClick={() => navigate(`/${tweet.user.display_id}/status/${tweet.id}`)}
    >
      <div className="TweetListItem__user__icon__wrapper">
        <UserIcon
          user={tweet.user}
          size={40}
          onClick={() => navigate(`/${tweet.user.display_id}`)}
        />
      </div>
      <div className="TweetListItem__content__wrapper">
        <div className="TweetListItem__content__user__wrapper">
          <span className="TweetListItem__content__user__name">
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
          <div className="TweetListItem__content__interactions__reply">
            <FaRegComment />
            <span>{tweet.interactions.reply_count}</span>
          </div>
          <div className="TweetListItem__content__interactions__retweet">
            <AiOutlineRetweet style={{ scale: "1.2" }} />
            <span>{tweet.interactions.retweet_count}</span>
          </div>
          <div
            className="TweetListItem__content__interactions__like"
            onClick={likeTweet}
          >
            <BsHeart />
            <span>{tweet.interactions.like_count}</span>
          </div>
        </div>
      </div>
    </div>
  );
};
