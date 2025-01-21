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
      <h1>タイムライン</h1>
      <p>ようこそ、{user.name}さん。</p>
      <button
        onClick={() =>
          navigate("/compose/tweet", { state: { background: location } })
        }
      >
        ツイートする
      </button>
      <TweetList tweets={tweets} refetch={fetchTweets} />
    </div>
  );
};
