import "./modal.css";
import { useNavigate } from "react-router";

export const Modal = ({ children }) => { 
  const navigate = useNavigate();
  return (
    <div className="modal-overlay" onClick={() => navigate(-1)}>
      <div className="modal-content" onClick={(e) => e.stopPropagation()}>
        {children}
      </div>
    </div>
  );
}
