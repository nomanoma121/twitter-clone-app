import { useEffect, useState } from "react";
import { serverFetch } from "../../../../utils/fetch";
import { SlCalender } from "react-icons/sl";
import { UserIcon } from "../../../../components/user-icon";
import { Button } from "../../../../components/button";
import "./user-profile.css";

type User = {
  id: number;
  display_id: string;
  name: string;
  icon_url: string;
  header_url: string;
  profile: string;
  follower_counts: number;
  followee_counts: number;
  created_at: string;
};

export const UserProfile = () => {
  const [user, setUser] = useState<User>({
    id: 0,
    display_id: "",
    name: "",
    icon_url: "",
    header_url: "",
    profile: "",
    follower_counts: 0,
    followee_counts: 0,
    created_at: "",
  });
  const path = window.location.pathname;
  const displayID = path.split("/")[1];
  const fetchUser = async () => {
    const res = await serverFetch(`/api/users/${displayID}`);
    if (res.ok) {
      const data = await res.json();
      setUser(data);
    }
  };

  useEffect(() => {
    fetchUser();
  }, []);

  const formatDate = (date: string) => {
    const d = new Date(date);
    return `${d.getFullYear()}年${d.getMonth()+1}月`;
  }

  console.log(user);

  return (
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
          <Button>フォロー</Button>
        </div>
        <div className="User__profile__name">{user.name}</div>
        <div className="User__profile__display-id">@{user.display_id}</div>
        <div className="User__profile__profile">{user.profile}</div>
        <div className="User__profile__created-at">
          <SlCalender />
          {" "+formatDate(user.created_at)}からTwitterを利用しています
        </div>
        <div className="User__profile__counts">
          <div>{user.followee_counts}<span>フォロー中</span></div>
          <div>{user.follower_counts}<span>フォロワー</span></div>
        </div>
      </div>
    </div>
  );
};
