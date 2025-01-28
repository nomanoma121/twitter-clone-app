import { ComponentProps } from "react";
import "./button.css";

type Props = ComponentProps<"button"> & {
  height?: number;
  width?: number;
  backgroundColor?: string;
  color?: string;
  active?: boolean;
};

export const Button = (props: Props) => {
  return (
    <button
      {...props}
      className={`Button ${props.active ? "active" : ""} ${props.className || ""}`}
      style={{
        height: props.height,
        width: props.width,
        backgroundColor: props.backgroundColor,
        color: props.color,
      }}
    />
  );
};
