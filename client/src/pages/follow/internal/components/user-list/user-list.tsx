import { TUser } from "../../../../../types/user";
import { UserListItem } from "../user-list-item/user-list-item";

type UserListProps = {
  users: TUser[];
  refetch: () => void;
};

export const UserList = ({
  users,
  refetch,
}: UserListProps) => {
  return (
    <div>
      {users.map((user) => (
        <UserListItem
          user={user}
          refetch={refetch}
          key={user.id}
        />
      ))}
    </div>
  );
};
