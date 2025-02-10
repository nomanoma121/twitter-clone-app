import { Button } from "../../../../../components/button";
import { UserIcon } from "../../../../../components/user-icon";
import { TUser } from "../../../../../types/user";
import { followAPI } from "../../../../../utils/followAPI";
import "./user-list-item.css";

type UserListProps = {
  user: TUser;
  refetch: () => void;
};

export const UserListItem = ({ user, refetch }: UserListProps) => {
  const handleFollow = async () => {
    const method = user?.followed_by_user ? "DELETE" : "POST";
    const res = await followAPI(user.display_id, method); 
    console.log(res);
    if (res.ok) {
      console.log("ok");
      refetch();
    }
  };

  return (
    <div className="UserListItem">
      <div className="UserListItem__icon">
        <UserIcon user={user} size={40} />
      </div>
      <div className="UserListItem__right">
        <div className="UserListItem__info">
          <div className="UserListItem__info__name">
            <div className="UserListItem__info__name__name">{user.name}</div>
            <div className="UserListItem__info__name__display-id">
              @{user.display_id}
            </div>
          </div>
          <div className="UserListItem__info__follow">
            <Button
              height={30}
              onClick={handleFollow}
              buttonActive={user.followed_by_user ? true : false}
            >
              {user.followed_by_user ? "フォロー中" : "フォロー"}
            </Button>
          </div>
        </div>
        <div className="UserListItem__profile">
          <div>{user.profile}</div>
        </div>
      </div>
    </div>
  );
};
