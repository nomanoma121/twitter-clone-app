export const User = () => {
  const path = window.location.pathname;
  const username = path.split("/")[1];
  console.log(username);
  return (
    <div>
      <div>
        <h1>{username}</h1>
        
      </div>
    </div>
  )
}
