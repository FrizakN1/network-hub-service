import React, {useEffect, useState} from "react";
import FetchRequest from "../fetchRequest";
import {FontAwesomeIcon} from "@fortawesome/react-fontawesome";
import {faCircle, faSquareCheck} from "@fortawesome/free-solid-svg-icons";
import NodesTable from "./NodesTable";

const ModalSelectTable = ({uri, type, selectRecord, setState, alreadySelect}) => {
    const [records, setRecords] = useState([])

    const handlerModalTableClose = (e) => {
        if (e.target.className === "modal-table") {
            setState(false)
        }
    }

    useEffect(() => {
        if (type !== "node") {
            FetchRequest("GET", uri, null)
                .then(response => {
                    if (response.success && response.data != null) {
                        setRecords(response.data)
                    }
                })
        }
    }, []);

    return (
        <div className="modal-table" onMouseDown={handlerModalTableClose}>
            <NodesTable action="select" selectFunction={(node) => {
                selectRecord(node)
                setState(false)
            }} selectNode={alreadySelect}/>
                {/*{(type === "owner" || type === "node_type") ?*/}
                {/*    <div className="contain">*/}
                {/*        <table>*/}
                {/*            <thead>*/}
                {/*            <tr className={"row-type-2"}>*/}
                {/*                <th>ID</th>*/}
                {/*                <th>Наименование</th>*/}
                {/*                <th>Дата создания</th>*/}
                {/*                <th></th>*/}
                {/*            </tr>*/}
                {/*            </thead>*/}
                {/*            <tbody>*/}
                {/*            {records.map((record, index) => (<tr key={index} className={index % 2 === 0 ? 'row-type-1' : 'row-type-2'}>*/}
                {/*                <td>{record.ID}</td>*/}
                {/*                <td>{record.Name}</td>*/}
                {/*                <td>{new Date(record.CreatedAt * 1000).toLocaleString().slice(0, 17)}</td>*/}
                {/*                <td>*/}
                {/*                    <FontAwesomeIcon icon={faCircle} title="Выбрать" onClick={() => {*/}
                {/*                        selectRecord(record)*/}
                {/*                        setState(false)*/}
                {/*                    }}/>*/}
                {/*                </td>*/}
                {/*            </tr>))}*/}
                {/*         */}
                {/*            </tbody>*/}
                {/*        </table>*/}
                {/*    </div>*/}
                {/*    :*/}
                {/*    <NodesTable action="select" selectFunction={(node) => {*/}
                {/*        selectRecord(node)*/}
                {/*        setState(false)*/}
                {/*    }}/>*/}
                {/*}*/}
        </div>
    )
}

export default ModalSelectTable