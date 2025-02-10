import { Container } from "../components/container";
import { Sidebar } from "../components/sidebar";
import { Outlet } from "react-router";

export const DefaultLayout = () => {
  return (
    <Container>
      <Sidebar />
      <div
        style={{
          maxWidth: "632px",
          width: "100%",
          borderLeft: "1px solid #e1e8ed",
          borderRight: "1px solid #e1e8ed",
        }}
      >
        <Outlet />
      </div>
    </Container>
  );
};
