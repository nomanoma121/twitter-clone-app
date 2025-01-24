import { TTweet } from "../../types";
import { BsHeart } from "react-icons/bs";
import { BsHeartFill } from "react-icons/bs";
import { serverFetch } from "../../utils/fetch";
import "./like-button.css";

type LikeButtonProps = {
  tweet: TTweet;
  refetch: () => void;
}

export const LikeButton = ({ tweet, refetch }: LikeButtonProps) => {
  const handleLike = async (e) => {
    e.stopPropagation();
    const endpoint = tweet.liked_by_user ? `/api/like/${tweet.id}` : `/api/like/${tweet.id}`;
    const method = tweet.liked_by_user ? "DELETE" : "POST";
    const res = await serverFetch(endpoint, {
      method,
      body: JSON.stringify({
        tweet_id: tweet.id,
      }),
      headers: {
        "Content-Type": "application/json",
      },
    });

    if (res.ok) {
      console.log("いいねしました");
      refetch();
    }
  };

  return (
    <div
      onClick={(e) => handleLike(e)}
      style={{ cursor: "pointer", zIndex: 999 }}
      className="LikeButton"
    >
      <div className="LikeIcon">
        {tweet.liked_by_user ? (
          <BsHeartFill color="rgb(249, 24, 128)" />
        ) : (
          <BsHeart />
        )}
      </div>
      <span>{tweet.interactions.like_count}</span>
    </div>
  );
};
