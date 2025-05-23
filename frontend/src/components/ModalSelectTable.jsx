import React from "react";
import NodesTable from "./NodesTable";

const ModalSelectTable = ({selectRecord, setState, alreadySelect}) => {
    const handlerModalTableClose = (e) => {
        if (e.target.className === "modal-table") {
            setState(false)
        }
    }

    return (
        <div className="modal-table" onMouseDown={handlerModalTableClose}>
            <NodesTable action="select" selectFunction={(node) => {
                selectRecord(node)
                setState(false)
            }} selectNode={alreadySelect}/>
        </div>
    )
}

export default ModalSelectTable