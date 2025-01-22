import "./modal.css";
export const Modal = ({ children }) => { 
  console.log("Modal");
  return (
    <div className="modal-overlay" onClick={() => window.history.back()}>
      <div className="modal-content" onClick={(e) => e.stopPropagation()}>
        {children}
      </div>
    </div>
  );
}
