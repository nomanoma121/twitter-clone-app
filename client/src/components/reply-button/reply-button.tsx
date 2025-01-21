import { TTweet } from "../../types";
import { FaRegComment } from "react-icons/fa";
import { useNavigate } from "react-router";
import { useLocation } from "react-router";

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
      style={{ cursor: "pointer", zIndex: 999 }}
    >
      <FaRegComment />
      <span>{tweet.interactions.reply_count}</span>
    </div>
  );
};
