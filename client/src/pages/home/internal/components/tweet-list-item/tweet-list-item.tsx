import { TTweet } from "../../../../../types/tweet";

interface TweetListItemProps {
  tweet: TTweet;
  refetch: () => void | null;
}

export const TweetListItem = ({ tweet, refetch }: TweetListItemProps) => {
  // refetchはとりあえず残したいけど今は使わない。エラーだけ外したい

  return (
    <div className="TweetListItem">
      <div className="TweetListItem__content">
        <p>{tweet.content}</p>
        <p className="TweetListItem__content__date">{tweet.createdAt}</p>
      </div>
    </div>
  );
};
