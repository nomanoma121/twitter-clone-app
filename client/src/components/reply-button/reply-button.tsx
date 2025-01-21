import { TTweet } from "../../types";
import { FaRegComment } from "react-icons/fa";

export const ReplyButton = ({ tweet }: TTweet) => {
  const handleReply = async (e) => {
    e.stopPropagation();
    console.log("リプライしました");
  }
  return (
    <div onClick={(e) => handleReply(e)} style={{ cursor: "pointer", zIndex: 999 }}>
      <FaRegComment />
      <span>{tweet.interactions.reply_count}</span>
    </div>
  );
};
