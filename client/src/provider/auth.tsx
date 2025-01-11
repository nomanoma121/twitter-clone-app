import React, {
  createContext,
  useCallback,
  useContext,
  useEffect,
  useState,
} from "react";
import { serverFetch } from "../utils/fetch";
import { useNavigate } from "react-router";

type User = {
  id: number;
  name: string;
  display_id: string;
  email: string;
};

type AuthContextType = {
  user: User | null;
  initialized: boolean;
  setUser: (user: User | null) => void;
  logout: () => void;
};

const AuthContext = createContext<AuthContextType | undefined>(undefined);

interface AuthProviderProps {
  children: React.ReactNode;
}

const fetchUser = async (): Promise<User | null> => {
  const res = await serverFetch("/auth/me");

  if (res.ok) {
    return res.json();
  }

  localStorage.removeItem("token");

  return null;
};

export const AuthProvider = ({ children }: AuthProviderProps) => {
  const [user, setUser] = useState<User | null>(null);
  const [initialized, setInitialized] = useState(false);
  const navigate = useNavigate();

  useEffect(() => {
    fetchUser().then((data) => {
      setUser(data);
      setInitialized(true);
    });
  }, []);

  const logout = useCallback(() => {
    localStorage.removeItem("token");
    setUser(null);
    navigate("/");
  }, [navigate]);

  return (
    <AuthContext.Provider
      value={{
        user,
        setUser,
        initialized,
        logout,
      }}
    >
      {children}
    </AuthContext.Provider>
  );
};

export const useAuth = () => {
  const context = useContext(AuthContext);
  if (context === undefined) {
    throw new Error("useAuth must be used within a AuthProvider");
  }
  return context;
};
