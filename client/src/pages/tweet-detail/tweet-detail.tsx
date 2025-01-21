import { TweetList } from "../../components/tweet-list";
import { useTweets } from "../../hooks/use-tweets";
import { useAuth } from "../../provider/auth";
import { useEffect, useState } from "react";
import { serverFetch } from "../../utils/fetch";
import { UserIcon } from "../../components/user-icon";
import RetweetItem from "../../components/retweet-item";
import { useNavigate } from "react-router";
import { TTweet } from "../../types/tweet";
import "./tweet-detail.css";

export const TweetDetail = () => {
  const [tweet, setTweet] = useState<TTweet | null>(null);
  const { user } = useAuth();
  const navigate = useNavigate();

  // 現在のパスからツイートIDを取得 /:username/status/:tweetId
  const path = window.location.pathname;
  const tweetId = path.split("/")[3];
  const { tweets, fetchTweets } = useTweets(`/api/tweets/${tweetId}/replies`);
  console.log(tweets);

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

  if (!user) return null;

  return (
    <div className="TweetDetail">
      <div className="TweetDetail__tweet">
        <div className="TweetDetail__user">
          <UserIcon
            user={tweet.user}
            size={40}
            onClick={() => navigate(`/${tweet.user.display_id}`)}
          />
          <div className="TweetDetail__user__info">
            <span className="TweetDetail__user__name">
              {tweet.user.name}
            </span>
            <span className="TweetDetail__user__displayID">
              @{tweet.user.display_id}
            </span>
          </div>
        </div>
        <div className="TweetDetail__content__wrapper">
          <div className="TweetDetail__content__content">{tweet.content}</div>
          {tweet.retweet && <RetweetItem retweet={tweet.retweet} />}
        </div>
      </div>
      <div>
        <TweetList tweets={tweets} refetch={fetchTweets} />
      </div>
    </div>
  );
};
