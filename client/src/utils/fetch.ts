export const serverFetch = (url: string, options: RequestInit = {}) => {
  let serverHost = import.meta.env.VITE_SERVER_HOST;
  if (!serverHost) {
    throw new Error("VITE_SERVER_HOST is not defined");
  }
  if (serverHost.endsWith("/")) {
    serverHost = serverHost.slice(0, -1);
  }
  const token = localStorage.getItem("token");
  return fetch(`${serverHost}${url}`, {
    ...options,
    headers: {
      ...(options.headers || {}),
      ...(token ? { Authorization: `Bearer ${token}` } : {}),
    },
  });
};
