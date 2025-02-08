import { Button } from "../../../../../components/button";
import { UserIcon } from "../../../../../components/user-icon";
import { TUser } from "../../../../../types/user";
import "./user-list-item.css";

type UserListProps = {
  user: TUser;
};

export const UserListItem = ({ user }: UserListProps) => {
  return (
    <div className="UserListItem">
      <div className="UserListItem__icon">
        <UserIcon user={user} size={40} />
      </div>
      <div className="UserListItem__right">
        <div  className="UserListItem__info">
          <div className="UserListItem__info__name">
            <div className="UserListItem__info__name__name">{user.name}</div>
            <div className="UserListItem__info__name__display-id">
              @{user.display_id}
            </div>
          </div>
          <div className="UserListItem__info__follow">
            <Button height={30} >フォロー</Button>
          </div>
        </div>
        <div className="UserListItem__profile">
          <div>{user.profile}</div>
        </div>
      </div>
    </div>
  );
};
