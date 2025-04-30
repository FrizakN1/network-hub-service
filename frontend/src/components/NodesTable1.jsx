import React, {useEffect, useState} from "react";
import UploadFile from "./UploadFile";
import {FontAwesomeIcon} from "@fortawesome/react-fontawesome";
import {
    faDownload,
    faEye,
    faFolderMinus,
    faFolderPlus,
    faPen,
    faPlus,
    faTrash
} from "@fortawesome/free-solid-svg-icons";
import NodeModalCreate from "./NodeModalCreate";
import {useNavigate, useParams} from "react-router-dom";
import FetchRequest from "../fetchRequest";

const NodesTable1 = () => {
    const [nodes, setNodes] = useState([])
    const [modalCreate, setModalCreate] = useState(false)
    const [modalEdit, setModalEdit] = useState({
        State: false,
        EditNode: {}
    })
    const { id } = useParams()
    const navigate = useNavigate()

    useEffect(() => {
        FetchRequest("GET", `/get_nodes/${id}`, null)
            .then(response => {
                if (response.success && response.data != null) {
                    setNodes(response.data)
                }
            })
    }, [id]);

    const handlerAddNode = (node) => {
        setNodes(prevState => [...prevState, node])
    }

    const handlerEditNode = (node) => {
        setNodes(prevState => prevState.map(_node => _node.ID === node.ID ? node : _node))
    }

    return (
        <div>
            {modalCreate && <NodeModalCreate action={"create"} setState={setModalCreate} returnNode={handlerAddNode}/>}
            {modalEdit.State && <NodeModalCreate action={"edit"} setState={(state) => setModalEdit(prevState => ({...prevState, State: state}))} editNode={modalEdit.EditNode} returnNode={handlerEditNode}/>}
            <div className="contain">
                <button className="add-node" onClick={() => setModalCreate(true)}>
                    <FontAwesomeIcon icon={faPlus}/> Добавить узел
                </button>
            </div>
            <div className="contain tables">
                {nodes.length > 0 ? (
                        <table className="nodes">
                            <thead>
                            <tr className={"row-type-2"}>
                                <th>Название узла</th>
                                <th>Дата создания</th>
                                <th></th>
                            </tr>
                            </thead>
                            <tbody>
                            {nodes.map((node, index) => (
                                <tr key={"node"+index} className={index % 2 === 0 ? 'row-type-1' : 'row-type-2'}>
                                    <td className={"col1"}>{node.Name}</td>
                                    <td className={"col2"}>{new Date(node.CreatedAt * 1000).toLocaleString().slice(0, 17)}</td>
                                    <td className={"col3"}>
                                        <FontAwesomeIcon icon={faEye} title="Просмотр" onClick={() => navigate(`/nodes/view/${node.ID}`)}/>
                                        <FontAwesomeIcon icon={faPen} title="Изменить" onClick={() => setModalEdit({State: true, EditNode: node})}/>
                                        <FontAwesomeIcon icon={faTrash} title="Удалить"/>
                                    </td>
                                </tr>
                            ))}
                            </tbody>
                        </table>
                    )
                    :
                    <div className="empty">Нет узлов</div>
                }
            </div>
        </div>
    )
}

export default NodesTable1