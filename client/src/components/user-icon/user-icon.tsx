import { useNavigate } from "react-router";
import "./user-icon.css";

type UserIconProps = {
  user: {
    id: number;
    name: string;
    display_id: string;
    icon_url: string;
  };
  size: number;
};

export const UserIcon = ({ user, size }: UserIconProps) => {
  const navigate = useNavigate();
  return (
    <img
      src={user.icon_url}
      className="UserIcon"
      alt="icon"
      height={size}
      width={size}
      onClick={(e) => {
        e.stopPropagation();
        navigate(`/${user.display_id}`, { state: { id: user.id } });
      }}
    />
  );
};
