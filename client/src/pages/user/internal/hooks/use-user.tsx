import { useEffect } from "react";
import { useState } from "react";
import { serverFetch } from "../../../../utils/fetch";

type User = {
  id: number;
  display_id: string;
  name: string;
  icon_url: string;
  header_url: string;
  profile: string;
  follower_counts: number;
  followee_counts: number;
  followed_by_user: boolean;
  created_at: string;
};

export const useUser = (displayID: string) => {
  const [userData, setUser] = useState<User | null>(null);

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

  return {
    userData,
    fetchUser,
  };
};
