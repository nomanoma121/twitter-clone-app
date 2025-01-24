import { useAuth } from "../../provider/auth";
import { TweetList } from "../../components/tweet-list";
import { useTweets } from "../../hooks/use-tweets";

export const Home = () => {
  const { user } = useAuth();
  const { tweets, fetchTweets } = useTweets("/api/tweets/timeline");

  if (!user) return null;

  return (
    <div className="Home">
      <TweetList tweets={tweets} refetch={fetchTweets} />
    </div>
  );
};
