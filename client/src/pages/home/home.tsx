import { useAuth } from "../../provider/auth";
import { TweetList } from "../../components/tweet-list";
import { useTweets } from "../../hooks/use-tweets";
import { useLocation, useNavigate } from "react-router";

export const Home = () => {
  const navigate = useNavigate();
  const location = useLocation();
  const { user } = useAuth();
  const { tweets, fetchTweets } = useTweets("/api/tweets/timeline");

  if (!user) return null;

  return (
    <div className="Home">
      <TweetList tweets={tweets} refetch={fetchTweets} />
    </div>
  );
};
