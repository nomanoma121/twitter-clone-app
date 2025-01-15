import { useEffect, useState } from "react";
import { serverFetch } from "../../../../utils/fetch";
import { Button } from "../../../../components/button";

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
  
  return (
        <div className="User" style={{ border: "1px solid black" }}>
          <div>
            <h1>{user.name}</h1>
            <h1>@{user.display_id}</h1>
            <p>{user.profile}</p>
            <Button>Follow</Button>
            <div>
              <p>Follower: {user.follower_counts}</p>
              <p>Following: {user.followee_counts}</p>
            </div>
            <div>
              <p>Joined {user.created_at}</p>
            </div>
          </div>
        </div>
  )
}
