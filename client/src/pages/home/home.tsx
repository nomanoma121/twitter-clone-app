import { useAuth } from "../../provider/auth";
import { TweetList } from "../../components/tweet-list";
import { useTweets } from "../../hooks/use-tweets";
import { useEffect } from "react";
import { useLocation } from "react-router";
import { TweetForm } from "../../components/tweet-form";

export const Home = () => {
  const { user } = useAuth();
  const { tweets, fetchTweets } = useTweets("/api/tweets/timeline");
  const location = useLocation();

  useEffect(() => {
    fetchTweets();
  }, [location]);

  if (!user) return null;

  return (
    <div className="Home">
      <TweetForm user={user} refetch={fetchTweets} />
      <TweetList tweets={tweets} refetch={fetchTweets} />
    </div>
  );
};
