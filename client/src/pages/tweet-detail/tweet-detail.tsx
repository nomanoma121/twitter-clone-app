import { TweetList } from "../../components/tweet-list";
import { useTweets } from "../../hooks/use-tweets";
import { useAuth } from "../../provider/auth";
import { useEffect, useState } from "react";
import { serverFetch } from "../../utils/fetch";
import { UserIcon } from "../../components/user-icon";
import RetweetItem from "../../components/retweet-item";
import { useNavigate } from "react-router";
import { TTweet } from "../../types/tweet";
import { Header } from "../../components/header";
import { LikeButton } from "../../components/like-button";
import { RetweetButton } from "../../components/retweet-button";
import { ReplyButton } from "../../components/reply-button";
import { Button } from "../../components/button";
import { useLocation } from "react-router";
import "./tweet-detail.css";

export const TweetDetail = () => {
  const [tweet, setTweet] = useState<TTweet | null>(null);
  const { user } = useAuth();
  const navigate = useNavigate();
  const location = useLocation();

  // 現在のパスからツイートIDを取得 /:username/status/:tweetId
  const path = location.pathname;
  const tweetId = path.split("/")[3];
  const { tweets, fetchTweets } = useTweets(`/api/tweets/${tweetId}/replies`);

  useEffect(() => {
    fetchTweets();
    fetchMainTweet();
  }, [location]);

  const fetchMainTweet = async () => {
    const res = await serverFetch(`/api/tweet/${tweetId}`);
    if (res.ok) {
      const data = await res.json();
      setTweet(data);
    }
  };

  useEffect(() => {
    fetchMainTweet();
  }, [tweetId]);

  if (!user) return null;

  // 午前6:50・2025年1月1日のような形式で返す関数
  const formatDate = (date: string) => {
    const d = new Date(date);
    // 午前か午後化の形式にする
    return `${d.getHours() > 12 ? "午後" : "午前"}${
      d.getHours() % 12
    }:${d.getMinutes()}・${d.getFullYear()}年${
      d.getMonth() + 1
    }月${d.getDate()}日`;
  };

  return (
    <>
    <Header title={"ツイートする"} />
    <div className="TweetDetail">
      {!tweet ? (
        <div>loading...</div>
      ) : (
        <>
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
              <div className="TweetDetail__content__content">
                {tweet.content}
              </div>
              {tweet.retweet && <RetweetItem retweet={tweet.retweet} />}
            </div>
            <div className="TweetDetail__date">
              <span>{formatDate(tweet.created_at)}</span>
            </div>
            <div className="TweetDetail__interactions__wrapper">
              <div className="TweetDetail__interactions">
                <ReplyButton
                  tweet={tweet}
                  className="TweetDetail__interactions__reply"
                />
                <RetweetButton
                  tweet={tweet}
                  className="TweetDetail__interactions__retweet"
                />
                <LikeButton
                  tweet={tweet}
                  refetch={fetchMainTweet}
                  className="TweetDetail__interactions__like"
                />
              </div>
            </div>
            <div>
              <UserIcon user={user} size={40} />
              <span>返信をツイート</span>
              <Button type="submit">リプライ</Button>
            </div>
          </div>
          <div>
            {!tweets ? (
              <div>loading...</div>
            ) : (
              <TweetList tweets={tweets} refetch={fetchTweets} />
            )}
          </div>
        </>
      )}
    </div>
  </>
  );
};
