import { useActionState, useCallback } from "react";
import { Input } from "../input";
import { Button } from "../button";
import { serverFetch } from "../../utils/fetch";

type TweetFormStateType = {
  message: string;
};

export const TweetForm = () => {
  const TweetAction = useCallback(
    async (
      _prevState: TweetFormStateType,
      formData: FormData
    ): Promise<TweetFormStateType> => {
      const content = formData.get("title");

      const res = await serverFetch("/api/tweet", {
        method: "POST",
        body: JSON.stringify({ content }),
        headers: {
          "Content-Type": "application/json",
        },
      });

      if (res.ok) {
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
    <div className="Tweet">
      <form action={submitAction} className="Tweet">
        <Input type="text" name="title" placeholder="What's happening?" />
        <Button type="submit">Tweet</Button>
      </form>
      {error && <p>{error.message}</p>}
    </div>
  );
};
