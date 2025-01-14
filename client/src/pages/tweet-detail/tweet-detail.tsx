import { TweetList } from "../../components/tweet-list";
import { useTweets } from "../../hooks/use-tweets";
import { useAuth } from "../../provider/auth";
import { TweetListItem } from "../../components/tweet-list-item";
import { useEffect, useState } from "react";
import { serverFetch } from "../../utils/fetch";

export const TweetDetail = () => {
  const [tweet, setTweet] = useState<Tweet | null>(null);
  const { user } = useAuth();

  // 現在のパスからツイートIDを取得 /:username/status/:tweetId
  const path = window.location.pathname;
  const tweetId = path.split("/")[3];
  const { tweets, fetchTweets } = useTweets(`/api/tweets/${tweetId}/replies`);

  useEffect(() => {
    const fetchTweet = async () => {
      const res = await serverFetch(`/api/tweet/${tweetId}`);
      if (res.ok) {
        const data = await res.json();
        setTweet(data);
      }
    };

    fetchTweet();
  }, [tweetId]);

  return (
    <div className="TweetDetail">
      <div className="TweetDetail__tweet">
        {tweet && <TweetListItem tweet={tweet} refetch={fetchTweets} />}
      </div>
      <div style={{ margin : "20px" }}>
        <TweetList tweets={tweets} refetch={fetchTweets} />
      </div>
    </div>
  );
};
