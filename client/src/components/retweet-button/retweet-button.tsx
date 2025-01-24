import { TTweet } from "../../types";
import { AiOutlineRetweet } from "react-icons/ai";
import { useNavigate } from "react-router";
import { useLocation } from "react-router";
import "./retweet-button.css";

export const RetweetButton = ({ tweet }: TTweet) => {
  const navigate = useNavigate();
  const location = useLocation();
  const handleRetweet = async (e) => {
    e.stopPropagation();
    navigate("/compose/tweet", {
      state: {
        background: location,
        tweet: {
          id: tweet.id,
          type: "retweet",
        },
      },
    });
  };

  return (
    <div
      onClick={(e) => handleRetweet(e)}
      style={{ cursor: "pointer", zIndex: 999 }}
      className="RetweetButton"
    >
      <div className="RetweetIcon">
        <AiOutlineRetweet style={{ scale: "1.2" }} />
      </div>
      <span>{tweet.interactions.retweet_count}</span>
    </div>
  );
};
