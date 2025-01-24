import { useEffect, useState } from "react";
import { serverFetch } from "../utils/fetch";

type Tweet = {
  id: number;
  user: {
    id: number;
    display_id: string;
    name: string;
    icon_url: string;
  };
  content: string;
  retweet: {
    id: number;
    user: {
      id: number;
      display_id: string;
      name: string;
      icon_url: string;
    };
    content: string | null;
    interactions: {
      retweet_count: number;
      reply_count: number;
      like_count: number;
    };
    created_at: string;
  } | null;
  interactions: {
    retweet_count: number;
    reply_count: number;
    like_count: number;
  };
  liked_by_user: boolean;
  created_at: string;
};

export const useTweets = (url: string) => {
  const [tweets, setTweets] = useState<Tweet[]>([]);

  const fetchTweets = async () => {
    const res = await serverFetch(url);

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
