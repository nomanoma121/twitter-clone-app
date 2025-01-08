import { useAuth } from "../../provider/auth";
import { TweetList } from "./internal/components/tweet-list";
import { useTweets } from "./internal/hook/use-tweets";

export const Home = () => {
  const { user } = useAuth();
  const { tweets, fetchTweets } = useTweets();

  if (!user) return null;

  return (
    <div className="Home">
      <h1>タイムライン</h1>
      <p>ようこそ、{user.name}さん。</p>
      <TweetList tweets={tweets} refetch={fetchTweets} />
    </div>
  );
}
