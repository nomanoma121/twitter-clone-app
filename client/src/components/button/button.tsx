import { ComponentProps } from "react";
import "./button.css";

type Props = ComponentProps<"button"> & {
  height?: number;
  width?: number;
  backgroundColor?: string;
  color?: string;
};

export const Button = (props: Props) => {
  return (
    <button
      {...props}
      className={["Button", props.className].join(" ")}
      style={{
        height: props.height,
        width: props.width,
        backgroundColor: props.backgroundColor,
        color: props.color,
      }}
    />
  );
};
