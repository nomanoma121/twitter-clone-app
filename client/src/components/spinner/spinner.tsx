import "./spinner.css";

export const Spinner = () => {
  return (
    <svg className="spinner" viewBox="0 0 50 50">
      <circle className="track" cx="25" cy="25" r="20" fill="none" strokeWidth="5" />
      <circle className="path" cx="25" cy="25" r="20" fill="none" strokeWidth="5" />
    </svg>
  );
};
