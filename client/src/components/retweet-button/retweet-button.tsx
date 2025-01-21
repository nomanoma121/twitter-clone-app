import React from "react";
import { TTweet } from "../../types";
import { AiOutlineRetweet } from "react-icons/ai";

export const RetweetButton = ({ tweet }: TTweet) => {
  const handleRetweet = async (e) => {
    e.stopPropagation();
    console.log("リツイートしました");
  }
  return (
    <div onClick={(e) => handleRetweet(e)} style={{ cursor: "pointer", zIndex: 999 }}>
      <AiOutlineRetweet style={{ scale: "1.2" }} />
      <span>{tweet.interactions.retweet_count}</span>
    </div>
  );
};
