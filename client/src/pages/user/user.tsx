import { useAuth } from "../../provider/auth";
import { UserProfile } from "./internal/components/user-profile";
import { UserTweets } from "./internal/components/user-tweets";
import { serverFetch } from "../../utils/fetch";
import { Header } from "../../components/header";
import { useEffect, useState } from "react";
import { useParams } from "react-router";
import{ useUser } from "./internal/hooks/use-user";

export const User = () => {
  const { user } = useAuth();
  const { displayID } = useParams();
  const [tweetCounts, setTweetCounts] = useState<number>(0);
  const { userData, fetchUser } = useUser(displayID);

  const fetchTweetCounts = async () => {
    const res = await serverFetch(`/api/users/${displayID}/tweet-counts`);
    if (res.ok) {
      const data = await res.json();
      setTweetCounts(data.tweet_counts);
    }
  };

  useEffect(() => {
    fetchTweetCounts();
  }, [user]);

  if (!user) return null;

  return (
    <div className="User">
      <Header title={userData?.name} subtitle={`${tweetCounts}件のツイート`} />
      <UserProfile user={userData} refetch={fetchUser} />
      <UserTweets />
    </div>
  );
};
