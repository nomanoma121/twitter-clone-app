import { ComponentProps } from "react";
import "./button.css";

type Props = ComponentProps<"button">;

export const Button = (props: Props) => {
  return <button {...props} className={["Button", props.className].join(" ")} />;
};
