import { serverFetch } from "../../../../utils/fetch";
import { SlCalender } from "react-icons/sl";
import { UserIcon } from "../../../../components/user-icon";
import { Button } from "../../../../components/button";
import "./user-profile.css";
import { useNavigate } from "react-router";

// TODO: Userの型を何とかする
type UserProfileProps = {
  user: any;
  refetch: () => void;
};

export const UserProfile = ({ user, refetch }: UserProfileProps) => {
  const navigate = useNavigate();
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
              <div className="User__profile__counts__item" onClick={() => navigate(`/${user.display_id}/following`)} >
                <div>{user.followee_counts}</div>
                <div className="User__profile__counts__follow">フォロー中</div>
              </div>
              <div className="User__profile__counts__item" onClick={() => navigate(`/${user.display_id}/followers`)} >
                <div>{user.follower_counts}</div>
                <div className="User__profile__counts__follow">フォロワー</div>
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
