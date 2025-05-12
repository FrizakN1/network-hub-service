import React, {useEffect, useState} from "react";
import FetchRequest from "../fetchRequest";
import SwitchModalCreate from "./SwitchModalCreate";
import {FontAwesomeIcon} from "@fortawesome/react-fontawesome";
import {faPen, faPlus} from "@fortawesome/free-solid-svg-icons";

const SwitchesPage = () => {
    const [switches, setSwitches] = useState([])
    const [modalCreate, setModalCreate] = useState(false)
    const [modalEdit, setModalEdit] = useState({
        State: false,
        EditSwitch: null
    })
    const [isLoaded, setIsLoaded] = useState(true)

    useEffect(() => {
        FetchRequest("GET", `/switches`, null)
            .then(response => {
                if (response.success) {
                    setSwitches(response.data != null ? response.data : [])
                }

                setIsLoaded(true)
            })
    }, []);

    const handlerAddSwitch = (record) => {
        setSwitches(prevState => [...prevState, record])
    }

    const handlerEditSwitch = (record) => {
        setSwitches(prevState => prevState.map(_record => record.ID === _record.ID ? record : _record))
    }

    return (
        <section className="references">
            {modalCreate && <SwitchModalCreate action="create" setState={setModalCreate} returnSwitch={handlerAddSwitch}/>}
            {modalEdit.State && <SwitchModalCreate action="edit"
                                                         setState={(state) => setModalEdit(prevState => ({...prevState, State: state}))}
                                                         returnSwitch={handlerEditSwitch}
                                                         editSwitch={modalEdit.EditSwitch}
            />}

            <div className="buttons">
                <button onClick={() => setModalCreate(true)}><FontAwesomeIcon icon={faPlus}/>
                    Создать модель коммутатора
                </button>
            </div>
            {isLoaded && <>{switches.length > 0 ? (
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
                        {switches.map((record, index) => (
                            <tr key={"record"+index} className={index % 2 === 0 ? 'row-type-1' : 'row-type-2'}>
                                <td>{record.ID}</td>
                                <td>{record.Name || record.TranslateValue}</td>
                                <td>{new Date(record.CreatedAt * 1000).toLocaleString().slice(0, 17)}</td>
                                <td>
                                    <FontAwesomeIcon icon={faPen} title="Редактировать" onClick={() => setModalEdit({State: true, EditSwitch: record})}/>
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

export default SwitchesPage