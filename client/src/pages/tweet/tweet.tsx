import { useActionState, useCallback } from "react";
import { Button } from "../../components/button";
import { serverFetch } from "../../utils/fetch";
import { useAuth } from "../../provider/auth";
import { useLocation } from "react-router";
import { useNavigate } from "react-router";
import { RxCross2 } from "react-icons/rx";
import { UserIcon } from "../../components/user-icon";
import "./tweet.css";

type TweetFormStateType = {
  message: string;
};

type TweetType = "tweet" | "reply" | "retweet";

export const Tweet = () => {
  const { user } = useAuth();
  const location = useLocation();
  const navigate = useNavigate();
  const tweetType: TweetType = location.state?.tweet.type;

  const TweetAction = useCallback(
    async (
      _prevState: TweetFormStateType,
      formData: FormData
    ): Promise<TweetFormStateType> => {
      const content = formData.get("content");

      const endpoint =
        tweetType === "tweet"
          ? "/api/tweet"
          : tweetType === "reply"
          ? `/api/tweet/${location.state?.tweet.id}/reply`
          : `/api/tweet/${location.state?.tweet.id}/retweet`;

      const res = await serverFetch(endpoint, {
        method: "POST",
        body: JSON.stringify({ content }),
        headers: {
          "Content-Type": "application/json",
        },
      });

      if (res.ok) {
        navigate(-1);
        return { message: "" };
      }

      return {
        message: "Tweetの投稿に失敗しました。",
      };
    },
    []
  );

  const [error, submitAction] = useActionState(TweetAction, {
    message: "",
  });

  if (!user) return null;

  return (
    <div className="Tweet">
      <div className="Tweet__top" >
        <div className="Tweet__top__cross" >
          <RxCross2 style={{ scale: "1.2"}} onClick={() => navigate(-1)} />
        </div>
      </div>
      <form action={submitAction} className="Tweet__form">
        <div className="Tweet__form__content" >
          <div className="Tweet__form__icon" >
            <UserIcon user={user} size={40} />
          </div>
          <div className="Tweet__form__input" >
            <input type="text" name="content" placeholder="いまどうしてる？" required />
          </div>
        </div>
        <div className="Tweet__form__footer" >
          <Button type="submit" buttonActive={true} >ツイート</Button>
        </div>
      </form>
      {error && <p>{error.message}</p>}
    </div>
  );
};
