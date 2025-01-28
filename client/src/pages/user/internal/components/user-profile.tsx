import { serverFetch } from "../../../../utils/fetch";
import { SlCalender } from "react-icons/sl";
import { UserIcon } from "../../../../components/user-icon";
import { Button } from "../../../../components/button";
import "./user-profile.css";

// TODO: Userの型を何とかする
type UserProfileProps = {
  user: any;
  refetch: () => void;
};

export const UserProfile = ({ user, refetch }: UserProfileProps) => {
  const handleFollow = async () => {
    const method = user?.followed_by_user ? "DELETE" : "POST";
    const endpoint = user?.followed_by_user ? "unfollow" : "follow";
    const res = await serverFetch(`/api/users/${user.display_id}/${endpoint}`, {
      method,
    });
    if (res.ok) {
      refetch();
    }
  };

  const formatDate = (date: string) => {
    const d = new Date(date);
    return `${d.getFullYear()}年${d.getMonth() + 1}月`;
  };

  return (
    <>
      {user ? (
        <div className="User__profile">
          <img
            src={user.header_url}
            alt="header"
            className="User__profile__header"
          />
          <div className="User__profile__icon">
            <UserIcon user={user} size={120} />
          </div>
          <div className="User__profile__info">
            <div className="User__profile__follow">
              <Button onClick={handleFollow} active={user?.followed_by_user}>
                {user?.followed_by_user ? "フォロー中" : "フォロー"}
              </Button>
            </div>
            <div className="User__profile__name">{user.name}</div>
            <div className="User__profile__display-id">@{user.display_id}</div>
            <div className="User__profile__profile">{user.profile}</div>
            <div className="User__profile__created-at">
              <SlCalender />
              {" " + formatDate(user.created_at)}からTwitterを利用しています
            </div>
            <div className="User__profile__counts">
              <div>
                {user.followee_counts}
                <span>フォロー中</span>
              </div>
              <div>
                {user.follower_counts}
                <span>フォロワー</span>
              </div>
            </div>
          </div>
        </div>
      ) : (
        <div>loading...</div>
      )}
    </>
  );
};
