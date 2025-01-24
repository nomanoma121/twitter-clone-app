import { TTweet } from "../../types";
import { BsHeart } from "react-icons/bs";
import { BsHeartFill } from "react-icons/bs";
import "./like-button.css";

export const LikeButton = ({ tweet }: TTweet) => {
  const handleLike = async (e) => {
    e.stopPropagation();
    console.log("いいねしました");
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
