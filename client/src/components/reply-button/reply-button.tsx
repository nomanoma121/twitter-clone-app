import { TTweet } from "../../types";
import { FaRegComment } from "react-icons/fa";
import { useNavigate } from "react-router";
import { useLocation } from "react-router";
import "./reply-button.css";

export const ReplyButton = ({ tweet }: TTweet) => {
  const navigate = useNavigate();
  const location = useLocation();
  const handleReply = async (e) => {
    e.stopPropagation();
    navigate("/compose/tweet", {
      state: {
        background: location,
        tweet: {
          id: tweet.id,
          type: "reply",
        },
      },
    });
  };
  return (
    <div
      onClick={(e) => handleReply(e)}
      className="ReplyButton"
    >
      <div className="ReplyIcon">
        <FaRegComment />
      </div>
      <span>{tweet.interactions.reply_count}</span>
    </div>
  );
};
