import { Button } from "../../components/button";
import { serverFetch } from "../../utils/fetch";

export const User = () => {
  const path = window.location.pathname;
  const username = path.split("/")[1];
  const res = await serverFetch(`/api/user/${username}`);
  return (
    <div>
      <div>
        <h1>{username}</h1>
        <Button>Follow</Button>
        <div>
          <p></p>
        </div>
      </div>
    </div>
  )
}
