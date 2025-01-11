import { TTweet } from "../../../../../types/tweet";

interface TweetListItemProps {
  tweet: TTweet;
  refetch: () => void | null; // エラーしないようにnullを追加
}

export const TweetListItem = ({ tweet, refetch }: TweetListItemProps) => {
  // refetchはとりあえず残したいけど今は使わない。エラーだけ外したかった
  return (
    <div className="TweetListItem" style={{ border: "1px solid black" }}>
      <p>{tweet.retweet ? "リツイート" : ""}</p>
      <div className="TweetListItem__user">
        <p>{tweet.user.name}</p>
        <p>@{tweet.user.display_id}</p>
      </div>
      <div className="TweetListItem__content">
        <p>{tweet.content}</p>
        <p className="TweetListItem__content__date">{tweet.created_at}</p>
      </div>
    </div>
  );
};
