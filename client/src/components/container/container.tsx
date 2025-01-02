import "./container.css";

interface Props {
  children: React.ReactNode;
}

export const Container = ({ children }: Props) => {
  return <div className="Container">{children}</div>;
};
