import { useNavigate } from "react-router";
import { TTweet } from "../../types/tweet";
import { serverFetch } from "../../utils/fetch";
import React from "react";

interface TweetListItemProps {
  tweet: TTweet;
  refetch: () => void; // エラーしないようにnullを追加
}

export const TweetListItem = ({ tweet, refetch }: TweetListItemProps) => {
  // refetchはとりあえず残したいけど今は使わない。エラーだけ外したかった
  const navigate = useNavigate();
  console.log(tweet.user.display_id);
  const likeTweet = async (e: React.MouseEvent) => {
    console.log(tweet.id);
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
  }

  return (
    <div key={tweet.id} className="TweetListItem" style={{ border: "1px solid black" }} onClick={() => navigate(`/${tweet.user.display_id}/status/${tweet.id}`)}>
      <p>{tweet.retweet ? "リツイート" : ""}</p>
      <div className="TweetListItem__user">
        <p>{tweet.user.name}</p>
        <p>@{tweet.user.display_id}</p>
      </div>
      <div className="TweetListItem__content">
        <p>{tweet.content}</p>
      </div>
      {tweet.retweet && (
        <div className="TweetListItem__retweet" style={{ border: "1px solid blue", margin: "10px" }} onClick={(e) => {
          e.stopPropagation();
          navigate(`/${tweet.retweet.user.display_id}/status/${tweet.retweet.id}`);
          }}>
          <div className="TweetListItem__retweet__user">
            <p>{tweet.retweet.user.name}</p>
            <p>@{tweet.retweet.user.display_id}</p>
          </div>
          <div className="TweetListItem__retweet__content">
            <p>{tweet.retweet.content}</p>
            <p className="TweetListItem__retweet__content__date">{tweet.retweet.created_at}</p>
          </div>
          <div className="TweetListItem__retweet__interactions" style={{ display: "flex" }}>
            <p>RT: {tweet.retweet.interactions.retweet_count}</p>
            <p>返信: {tweet.retweet.interactions.reply_count}</p>
            <p>いいね: {tweet.retweet.interactions.like_count}</p>
          </div>
          <p>{tweet.retweet.created_at}</p>
        </div>
      )}
      <div className="TweetListItem__interactions" style={{ display: "flex" }}>
        <p>RT: {tweet.interactions.retweet_count}</p>
        <p>返信: {tweet.interactions.reply_count}</p>
        <div onClick= {(e) => likeTweet(e)} style={{ backgroundColor: "pink", zIndex: "9999999"}}>いいね: {tweet.interactions.like_count}</div>
      </div>
      <p>{tweet.created_at}</p>
    </div>
  );
};
