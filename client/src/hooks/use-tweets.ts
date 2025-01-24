import { useState } from "react";
import { serverFetch } from "../utils/fetch";
import { TTweet } from "../types/tweet";

export const useTweets = (url: string) => {
  const [tweets, setTweets] = useState<TTweet[]>([]);

  const fetchTweets = async () => {
    const res = await serverFetch(url);

    if (res.ok) {
      const data = await res.json();
      setTweets(data);
    }
  };

  return {
    tweets,
    fetchTweets,
  };
};
