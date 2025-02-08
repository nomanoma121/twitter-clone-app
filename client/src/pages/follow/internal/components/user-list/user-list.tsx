import { TUser } from "../../../../../types/user";
import { UserListItem } from "../user-list-item/user-list-item";

type UserListProps = {
  users: TUser[];
};
export const UserList = ( {users}: UserListProps) => {
  return (
    <div>
      {users.map((user) => (
        <UserListItem user={user} key={user.id} />
      ))}
    </div>
  )
}
