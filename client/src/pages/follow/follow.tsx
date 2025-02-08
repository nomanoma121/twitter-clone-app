import { Tabbar } from "../../components/tabbar/tabbar";
import { UserList } from "./internal/components/user-list/user-list";
import { useNavigate, useParams } from "react-router";
import { serverFetch } from "../../utils/fetch";
import { TUser } from "../../types/user";
import { useEffect, useState } from "react";

type FollowProps = {
  viewMode: "following" | "followers";
};

export const Follow = ({ viewMode }: FollowProps) => {
  const { displayID } = useParams();
  const navigate = useNavigate();
  const [follows, setFollows] = useState<TUser | null>(null);
  if (!displayID) {
    return <div>404 Not Found</div>;
  }

  const switchTab = () => {
    navigate(
      `/${displayID}/${viewMode === "following" ? "followers" : "following"}`
    );
  };

  const fetchFollows = async () => {
    const res = await serverFetch(
      `/api/users/${displayID}/${
        viewMode === "following" ? "followees" : "followers"
      }`
    );
    if (res.ok) {
      const data = await res.json();
      setFollows(data);
      console.log(data);
    }
  };

  useEffect(() => {
    fetchFollows();
  }, [viewMode]);

  if (!follows) {
    return <div>loading...</div>;
  }

  return (
    <div>
      <Tabbar
        titles={{ first: "フォロワー", second: "フォロー中" }}
        switchTab={switchTab}
        defaultTab={viewMode === "following" ? "second" : "first"}
      />
      <UserList users={follows} />
    </div>
  );
};
