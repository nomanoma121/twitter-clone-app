import { useAuth } from "../../provider/auth";

export const TweetDetail = () => {
  const { user } = useAuth();
  
  return (
    <div className="TweetDetail">
      <h1>ツイート詳細</h1>
    </div>
  );
}
