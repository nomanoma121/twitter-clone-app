import { useNavigate } from "react-router";
import { UserIcon } from "../user-icon";
import { displayTime } from "../../utils/display-time";
import "./retweet-item.css";

type RetweetItemProps = {
  retweet: {
    id: number;
    user: {
      id: number;
      name: string;
      display_id: string;
      icon_url: string;
    };
    content: string;
    interactions: {
      retweet_count: number;
      reply_count: number;
      like_count: number;
    };
    created_at: string;
  };
};

export const RetweetItem = ({ retweet }: RetweetItemProps) => {
  const navigate = useNavigate();
  return (
    <div
      className="RetweetItem"
      onClick={(e) => {
        navigate(`/retweet.display_id/status/${retweet.id}`);
        e.stopPropagation();
      }}
    >
      <div className="RetweetItem__user__wrapper">
        <UserIcon
          user={retweet.user}
          size={24}
          onClick={() => navigate(`/${retweet.user.display_id}`)}
        />
        <span className="RetweetItem__user__name">{retweet.user.name}</span>
        <span className="RetweetItem__user__displayID">
          @{retweet.user.display_id}
        </span>
        <span className="RetweetItem__user__createdAt">
          ・{displayTime(retweet.created_at)}
        </span>
      </div>
      <div className="RetweetItem__content">{retweet.content}</div>
    </div>
  );
};
