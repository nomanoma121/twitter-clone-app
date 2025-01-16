import { useNavigate } from "react-router"
import UserIcon from "../user-icon"
import { displayTime } from "../../utils/display-time"
import { FaRegComment } from "react-icons/fa6"
import { AiOutlineRetweet } from "react-icons/ai"
import { BsHeart } from "react-icons/bs"
import "./retweet-item.css"

type RetweetItemProps = {
  retweet: {
    id: number
    user: {
      id: number
      name: string
      display_id: string
      icon_url: string
    }
    content: string
    interactions: {
      retweet_count: number
      reply_count: number
      like_count: number
    }
    created_at: string
  }
}

export const RetweetItem = ({ retweet }: RetweetItemProps) => {
  const navigate = useNavigate();
  return (
    <div className="RetweetItem">
      <div className="RetweetItem__user__icon__wrapper">
        <UserIcon
          user={retweet.user}
          size={24}
          onClick={() => navigate(`/${retweet.user.display_id}`)}
        />
      </div>
      <div className="RetweetItem__content__wrapper">
        <div className="RetweetItem__content__user__wrapper">
          <span className="RetweetItem__content__user__name">
            {retweet.user.name}
          </span>
          <span className="RetweetItem__content__user__displayID">
            @{retweet.user.display_id}
          </span>
          <span className="RetweetItem__content__user__createdAt">
            ・{displayTime(retweet.created_at)}
          </span>
        </div>
        <div className="RetweetItem__content__content">{retweet.content}</div>
        <div className="RetweetItem__content__interactions">
          <div className="RetweetItem__content__interactions__reply">
            <FaRegComment />
            <span>{retweet.interactions.reply_count}</span>
          </div>
          <div className="RetweetItem__content__interactions__retweet">
            <AiOutlineRetweet style={{ scale: "1.2" }} />
            <span>{retweet.interactions.retweet_count}</span>
          </div>
          <div
            className="RetweetItem__content__interactions__like"
            // onClick={likeTweet}
          >
            <BsHeart />
            <span>{retweet.interactions.like_count}</span>
          </div>
        </div>
      </div>
    </div>
  );
}
