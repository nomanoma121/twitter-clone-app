import { useActionState, useCallback } from "react";
import { Button } from "../../components/button";
import { Input } from "../../components/input";
import { serverFetch } from "../../utils/fetch";
import { useAuth } from "../../provider/auth";
import { useLocation } from "react-router";

type TweetFormStateType = {
  message: string;
};

type TweetType = "tweet" | "reply" | "retweet";

export const Tweet = () => {
  const { user } = useAuth();
  const location = useLocation();
  const tweetType: TweetType = location.state?.tweet.type;

  const TweetAction = useCallback(
    async (
      _prevState: TweetFormStateType,
      formData: FormData
    ): Promise<TweetFormStateType> => {
      const content = formData.get("title");

      let requestUrl;
      if (tweetType === "tweet") {
        requestUrl = "/api/tweet";
      } else if (tweetType === "reply") {
        requestUrl = `/api/tweet/${location.state?.tweet.id}/reply`;
      } else if (tweetType === "retweet") {
        requestUrl = `/api/tweet/${location.state?.tweet.id}/retweet`;
      } else {
        return {
          message: "Tweetの投稿に失敗しました。",
        };
      }

      const res = await serverFetch(requestUrl, {
        method: "POST",
        body: JSON.stringify({ content }),
        headers: {
          "Content-Type": "application/json",
        },
      });

      if (res.ok) {
        window.history.back();
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
      <form action={submitAction} className="Tweet">
        <Input type="text" name="title" placeholder="What's happening?" />
        <Button type="submit">Tweet</Button>
      </form>
      {error && <p>{error.message}</p>}
    </div>
  );
};
