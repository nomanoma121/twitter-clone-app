import { ComponentProps } from "react";
import "./input.css";

type Props = ComponentProps<"input">;

export const Input = (props: Props) => {
  return <input {...props} className={["Input", props.className].join(" ")} />;
};
