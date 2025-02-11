import { TTweet } from "../../types/tweet";
import { TweetListItem } from "../tweet-list-item";

interface TweetListProps {
  tweets: TTweet[];
  refetch: () => void;
}

export const TweetList = ({ tweets, refetch }: TweetListProps) => {
  return (
    <div>
      {tweets.slice().reverse().map((tweet) => (
        <TweetListItem key={`${tweet.user.display_id} + ${tweet.id}}`} tweet={tweet} refetch={refetch} />
      ))}
    </div>
  );
};
