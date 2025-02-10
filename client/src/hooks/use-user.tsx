import { useEffect } from "react";
import { useState } from "react";
import { serverFetch } from "../utils/fetch";
import { TUser } from "../types/user";

export const useUser = (displayID: string) => {
  const [userData, setUser] = useState<TUser | null>(null);

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
