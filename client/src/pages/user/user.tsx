import { useAuth } from "../../provider/auth";
import { UserProfile } from "./internal/components/user-profile";

export const User = () => {
  const { user } = useAuth();

  if (!user) return null;

  return (
    <div className="User" style={{ border: "1px solid black" }}>
      <UserProfile />
    </div>
  );
};
