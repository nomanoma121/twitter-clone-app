import { useEffect } from "react";
import { TweetList } from "../../../../components/tweet-list"
import { useTweets } from "../../../../hooks/use-tweets";
import { useLocation } from "react-router";


export const UserTweets = () => {
  const location = useLocation();
  const path = location.pathname
  const displayId = path.split("/")[1];
  const { tweets, fetchTweets } = useTweets(`/api/users/${displayId}/tweets`);

  useEffect(() => {
    fetchTweets();
  }, [location]);

  return (
    <div>
      <TweetList tweets={tweets} refetch={fetchTweets} />
    </div>
  )
}
