import { useParams } from "react-router";
import { serverFetch } from "../../../../utils/fetch";
import { SlCalender } from "react-icons/sl";
import { UserIcon } from "../../../../components/user-icon";
import { Button } from "../../../../components/button";
import { useUser } from "../hooks/use-user";
import "./user-profile.css";

export const UserProfile = () => {
  const { displayID } = useParams();
  if (!displayID) {
    return null;
  }
  const { user, fetchUser } = useUser(displayID);

  const handleFollow = async () => {
    const res = await serverFetch(`/api/users/${displayID}/follow`, {
      method: "POST",
    });
    if (res.ok) {
      fetchUser();
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
              <Button onClick={handleFollow}>フォロー</Button>
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
