import { serverFetch } from "./fetch";

export const followAPI = async (followedUserID: string, method: "POST" | "DELETE") => {
  const endpoint = method === "POST" ? "follow" : "unfollow";
  const res = await serverFetch(`/api/users/${followedUserID}/${endpoint}`, {
    method,
  });
  return res;
}
