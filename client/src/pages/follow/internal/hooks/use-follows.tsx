import { useEffect, useState } from 'react';
import { serverFetch } from '../../../../utils/fetch';
import { TUser } from '../../../../types/user';

type TUrl = "following" | "followers";

export const useFollows = (displayID: string, mode: TUrl) => {
  const [follows, setFollows] = useState<TUser[] | null>(null);
  const url = mode === "following" ? "followees" : "followers";
  
  const fetchFollows = async () => {
    const res = await serverFetch(`/api/users/${displayID}/${url}`);
    if (res.ok) {
      const data = await res.json();
      setFollows(data);
    }
  };

  useEffect(() => {
    fetchFollows();
  }, [displayID, mode]);

  return {
    follows,
    fetchFollows,
  };
}
