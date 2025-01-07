import { useEffect, useState } from "react";
import { serverFetch } from "../../../../utils/fetch";

type Tweet = {
  id: number
  user: {
    id: number
    name: string
    email: string
  }
  content: string
  retweet: Tweet | null
}

export const useTweets = () => {
  const [tweets, setTweets] = useState<Tweet[]>([]);

  const fetchTweets = async () => {
    const res = await serverFetch("/api/tweets");

    if (res.ok) {
      const data = await res.json();
      setTweets(data);
    }
  };

  useEffect(() => {
    fetchTweets();
  }, []);

  return {
    tweets,
    fetchTweets,
  };
};
