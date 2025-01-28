import { useEffect, useState } from "react";
import { TweetList } from "../../../../components/tweet-list"
import { useTweets } from "../../../../hooks/use-tweets";
import { useLocation } from "react-router";
import { Tabbar } from "../../../../components/tabbar";

export const UserTweets = () => {
  const location = useLocation();
  const displayId = location.pathname.split("/")[1];
  const [endpoint, setEndpoint] = useState(`/api/users/${displayId}/tweets`);
  const { tweets, fetchTweets } = useTweets(endpoint);

  useEffect(() => {
    fetchTweets();
  }, [location, endpoint]);

  const switchTab = () => {
    endpoint === `/api/users/${displayId}/tweets`
      ? setEndpoint(`/api/users/${displayId}/replies`)
      : setEndpoint(`/api/users/${displayId}/tweets`);
  }

  return (
    <div>
      <Tabbar titles={{first: "ツイート", second: "返信"}} switchTab={switchTab} />
      <TweetList tweets={tweets} refetch={fetchTweets} />
    </div>
  )
}
