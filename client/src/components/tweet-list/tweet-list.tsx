import { TTweet } from "../../types/tweet";
import { TweetListItem } from "../tweet-list-item";

interface TweetListProps {
  tweets: TTweet[];
  refetch: () => void;
}

export const TweetList = ({ tweets, refetch }: TweetListProps) => {
  return (
    <div>
      {tweets.map((tweet, index) => (
        <TweetListItem key={tweet.id * (index + 1)} tweet={tweet} refetch={refetch} />
      ))}
    </div>
  );
};
