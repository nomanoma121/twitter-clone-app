import { useActionState, useCallback } from "react";
import { Button } from "../../components/button";
import { Input } from "../../components/input";
import { serverFetch } from "../../utils/fetch";
import { useAuth } from "../../provider/auth";

type TweetFormStateType = {
  message: string;
};

export const Tweet = () => {
  const { user } = useAuth();

  const TweetAction = useCallback(
    async (
      _prevState: TweetFormStateType,
      formData: FormData
    ): Promise<TweetFormStateType> => {
      const content = formData.get("title");

      console.log(content);

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

  if (!user) return null;

  return (
    <div>
      <form action={submitAction} className="Tweet">
        <Input type="text" name="title" placeholder="What's happening?" />
        <Button type="submit">Tweet</Button>
      </form>
      {error && <p>{error.message}</p>}
    </div>
  );
};
