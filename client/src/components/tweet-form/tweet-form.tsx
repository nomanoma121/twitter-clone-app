import { useActionState, useCallback } from "react";
import { Button } from "../button";
import { serverFetch } from "../../utils/fetch";
import { UserIcon } from "../user-icon";
import "./tweet-form.css";

type TweetFormStateType = {
  message: string;
};

// type User = {
//   id: number;
//   name: string;
//   display_id: string;
//   icon_url: string;
// }

// TODO: userに型を付ける
export const TweetForm = ({ user, refetch }) => {
  const TweetAction = useCallback(
    async (
      _prevState: TweetFormStateType,
      formData: FormData
    ): Promise<TweetFormStateType> => {
      const content = formData.get("content");

      const res = await serverFetch("/api/tweet", {
        method: "POST",
        body: JSON.stringify({ content }),
        headers: {
          "Content-Type": "application/json",
        },
      });

      if (res.ok) {
        refetch();
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

  return (
    <div className="TweetForm">
      <div className="TweetForm__user">
        <UserIcon user={user} size={40} />
      </div>
      <form action={submitAction} className="TweetForm__form">
        <input type="text" name="content" className="TweetForm__form__input" placeholder="いまどうしてる？" />
        <div className="TweetForm__form__border"></div>
        <Button type="submit" className="TweetForm__form__button" buttonActive={true}>ツイートする</Button>
      </form>
      {error && <p>{error.message}</p>}
    </div>
  );
};
