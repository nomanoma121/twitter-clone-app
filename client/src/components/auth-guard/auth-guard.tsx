import { useNavigate } from "react-router";
import { useAuth } from "../../provider/auth";

interface Props {
  children: React.ReactNode;
  redirectPath?: string;
}

export const AuthGuard = ({ children, redirectPath = "/" }: Props) => {
  const auth = useAuth();
  const navitate = useNavigate();
  if (!auth.initialized) return null;
  if (!auth.user) {
    navitate(redirectPath);
    return null;
  }
  return <>{children}</>;
};
