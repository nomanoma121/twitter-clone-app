import { TweetList } from "../../../../components/tweet-list"
import { useTweets } from "../../../../hooks/use-tweets";


export const UserTweets = () => {
  const path = window.location.pathname;
  const displayID = path.split("/")[1];
  const { tweets, refetch } = useTweets(`/api/users/${displayID}/tweets`);
  return (
    <div style={{ margin: "20px" }}>
      <TweetList tweets={tweets} refetch={refetch} />
    </div>
  )
}
