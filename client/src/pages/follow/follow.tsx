import { Tabbar } from "../../components/tabbar/tabbar";
import { UserList } from "./internal/components/user-list/user-list";
import { useNavigate, useParams, useLocation } from "react-router";
import { useFollows } from "./internal/hooks/use-follows";

export const Follow = () => {
  const { displayID } = useParams();
  const location = useLocation();
  const navigate = useNavigate();
  const currentPath = location.pathname.split("/").pop();

  if (!displayID) {
    return <div>loading...</div>;
  }

  console.log(displayID, currentPath);
  const { follows, fetchFollows } = useFollows(displayID, currentPath as "following" | "followers");

  const switchTab = () => {
    navigate(
      `/${displayID}/${currentPath === "followers" ? "following" : "followers"}`
    );
  };

  if (!follows) {
    return <div>loading...</div>;
  }

  console.log(follows);

  return (
    <div>
      <Tabbar
        titles={{ first: "フォロワー", second: "フォロー中" }}
        switchTab={switchTab}
        defaultTab={currentPath === "following" ? "second" : "first"}
      />
      <UserList users={follows} refetch={fetchFollows} />
    </div>
  );
};
