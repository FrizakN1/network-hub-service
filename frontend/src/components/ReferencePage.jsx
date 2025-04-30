import React, {useEffect, useState} from "react";
import {FontAwesomeIcon} from "@fortawesome/react-fontawesome";
import {faPen, faPlus} from "@fortawesome/free-solid-svg-icons";
import FetchRequest from "../fetchRequest";
import NodeReferenceRecordModalCreate from "./NodeReferenceRecordModalCreate";

const ReferencePage = ({reference}) => {
    const [records, setRecords] = useState([])
    const [modalCreate, setModalCreate] = useState(false)
    const [modalEdit, setModalEdit]= useState({
        State: false,
        EditRecord: null
    })
    const [isLoaded, setIsLoaded] = useState(true)

    useEffect(() => {
        FetchRequest("GET", `/get_${reference}`, null)
            .then(response => {
                if (response.success && response.data != null) {
                    setRecords(response.data)
                }

                setIsLoaded(true)
            })
    }, []);

    const handlerAddRecord = (record) => {
        setRecords(prevState => [...prevState, record])
    }

    const handlerEditRecord = (record) => {
        setRecords(prevState => prevState.map(_record => record.ID === _record.ID ? record : _record))
    }

    return (
        <section className="references">
            {modalCreate && <NodeReferenceRecordModalCreate action="create" setState={setModalCreate} returnRecord={handlerAddRecord} reference={reference}/>}
            {modalEdit.State && <NodeReferenceRecordModalCreate action="edit"
                                                 setState={(state) => setModalEdit(prevState => ({...prevState, State: state}))}
                                                 returnRecord={handlerEditRecord}
                                                 editRecord={modalEdit.EditRecord}
                                                 reference={reference}
            />}
            <div className="buttons">
                <button onClick={() => setModalCreate(true)}><FontAwesomeIcon icon={faPlus}/> {reference === "node_types" ? "Создать тип узла" : "Создать владельца"}</button>
            </div>
            {isLoaded && <>{records.length > 0 ? (
                    <table>
                        <thead>
                        <tr className={"row-type-2"}>
                            <th>ID</th>
                            <th>Наименование</th>
                            <th>Дата создания</th>
                            <th></th>
                        </tr>
                        </thead>
                        <tbody>
                        {records.map((record, index) => (
                            <tr key={"record"+index} className={index % 2 === 0 ? 'row-type-1' : 'row-type-2'}>
                                <td>{record.ID}</td>
                                <td>{record.Name}</td>
                                <td>{new Date(record.CreatedAt * 1000).toLocaleString().slice(0, 17)}</td>
                                <td>
                                    <FontAwesomeIcon icon={faPen} title="Редактировать" onClick={() => setModalEdit({State: true, EditRecord: record})}/>
                                </td>
                            </tr>
                        ))}
                        </tbody>
                    </table>
                )
                :
                <div className="empty">Таблица пуста</div>
            }</>}
        </section>
    )
}

export default ReferencePage