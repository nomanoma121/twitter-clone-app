import { useAuth } from "../../provider/auth";
import { TweetList } from "../../components/tweet-list";
import { useTweets } from "../../hooks/use-tweets";
import { useEffect, useState } from "react";
import { useLocation } from "react-router";
import { TweetForm } from "../../components/tweet-form";
import { Tabbar } from "../../components/tabbar";

export const Home = () => {
  const { user } = useAuth();
  const [endpoint, setEndpoint] = useState("/api/tweets/timeline");
  const { tweets, fetchTweets } = useTweets(endpoint);
  const location = useLocation();

  const switchTab = () => {
    endpoint === "/api/tweets/timeline"
      ? setEndpoint("/api/tweets/follow")
      : setEndpoint("/api/tweets/timeline");
  };

  useEffect(() => {
    fetchTweets();
  }, [location, endpoint]);

  if (!user) return null;

  return (
    <div className="Home">
      <Tabbar
        titles={{
          first: "タイムライン",
          second: "フォロー中",
        }}
        switchTab={switchTab}
      />
      <TweetForm user={user} refetch={fetchTweets} />
      <TweetList tweets={tweets} refetch={fetchTweets} />
    </div>
  );
};
