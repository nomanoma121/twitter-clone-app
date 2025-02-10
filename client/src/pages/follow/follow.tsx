import { Tabbar } from "../../components/tabbar/tabbar";
import { Header } from "../../components/header/header";
import { UserList } from "./internal/components/user-list/user-list";
import { useNavigate, useParams, useLocation } from "react-router";
import { useFollows } from "./internal/hooks/use-follows";
import { useUser } from "../../hooks/use-user";

export const Follow = () => {
  const { displayID } = useParams();
  const location = useLocation();
  const navigate = useNavigate();
  const currentPath = location.pathname.split("/").pop();

  if (!displayID) {
    return <div>loading...</div>;
  }

  const { follows, fetchFollows } = useFollows(displayID, currentPath as "following" | "followers");
  const { userData } = useUser(displayID);

  const switchTab = () => {
    navigate(
      `/${displayID}/${currentPath === "followers" ? "following" : "followers"}`
    );
  };

  if (!follows || !userData) {
    return <div>loading...</div>;
  }

  return (
    <div>
      <Header title={userData.name} subtitle={`@${userData.display_id}`} />
      <Tabbar
        titles={{ first: "フォロワー", second: "フォロー中" }}
        switchTab={switchTab}
        defaultTab={currentPath === "following" ? "second" : "first"}
      />
      <UserList users={follows} refetch={fetchFollows} />
    </div>
  );
};
