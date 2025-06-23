import React, {useContext, useEffect, useRef, useState} from "react";
import {FontAwesomeIcon} from "@fortawesome/react-fontawesome";
import {
    faAngleLeft,
    faAngleRight,
    faAnglesLeft, faAnglesRight,
    faEye,
    faPen, faPlus, faRotate, faSquareCheck,
    faTrash
} from "@fortawesome/free-solid-svg-icons";
import {
    faCircle, faCircleDot
} from "@fortawesome/free-regular-svg-icons"
import FetchRequest from "../fetchRequest";
import NodeModalCreate from "./NodeModalCreate";
import {useNavigate} from "react-router-dom";
import AuthContext from "../context/AuthContext";

const NodesTable = ({houseID = 0, canCreate = false, action, selectFunction, selectNode, defaultAddress}) => {
    const searchDebounceTimer = useRef(0)
    const [nodes, setNodes] = useState([])
    const [isLoaded, setIsLoaded] = useState(false)
    const [currentPage, setCurrentPage] = useState(1)
    const [allPage, setAllPage] = useState([])
    const [showPages, setShowPages] = useState(null)
    const [search, setSearch] = useState("")
    const [count, setCount] = useState(0)
    const [modalCreate, setModalCreate] = useState(false)
    const [modalEdit, setModalEdit] = useState({
        State: false,
        EditNodeID: 0
    })
    const navigate = useNavigate()
    const { user } = useContext(AuthContext)

    useEffect(() => {
        handlerGetNodes()
    }, []);

    const handlerGetNodes = (value = "") => {
        let uri = "/nodes"
        let params = new URLSearchParams({
            offset: String((currentPage-1)*20),
            only_active: String(action === "select"),
        })

        if (value.length > 0) {
            params.set("search", value)
            uri = "/nodes/search"
        }

        if (houseID > 0) {
            uri = `/houses/${houseID}/nodes`
        }

        FetchRequest("GET", `${uri}?${params.toString()}`, null)
            .then(response => {
                if (response.success) {
                    setNodes(response.data.Nodes != null ? response.data.Nodes : [])
                    setCount(response.data.Count)
                    setIsLoaded(true)
                }
            })
    }

    const handlerSearch = (e) => {
        setSearch(e.target.value)

        clearTimeout(searchDebounceTimer.current)

        searchDebounceTimer.current = setTimeout(() =>  {
            handlerGetNodes(e.target.value)
        }, 500)
    }

    useEffect(() => {
        setAllPage(Array.from({ length: Math.ceil(count/20) }, (_, index) => (
            <span
                key={index}
                className={index + 1 === currentPage ? 'active' : ''}
                onClick={() => setCurrentPage(index + 1)}
            >
                {index + 1}
            </span>
        )))

        if (Math.ceil(count/20) !== 0 && Math.ceil(count/20) < currentPage) {
            setCurrentPage(Math.ceil(count/20))
        }
    }, [count, currentPage]);

    useEffect(() => {
        handlerGetNodes(search)

        if (allPage.length <= 7) {
            setShowPages(
                <div className="pages">
                    {allPage}
                </div>)
        } else {
            let pagesSlice = [allPage[0]]

            if (currentPage > 4) {
                pagesSlice.push(<span key={"space-1"}>...</span>)
            }

            let startPoint = currentPage - 3
            let endPoint = currentPage + 2

            if (startPoint <= 1) {
                endPoint += (startPoint-1) * -1
                startPoint = 1
            } else if (endPoint >= allPage.length) {
                startPoint -= (allPage.length-currentPage-3) * -1
                endPoint = allPage.length-1
            }

            for (let i = startPoint; i < endPoint; i++) {
                pagesSlice.push(allPage[i])
            }

            if (currentPage <= allPage.length-4) {
                pagesSlice.push(<span key={"space-2"}>...</span>)
            }

            pagesSlice.push(allPage[allPage.length-1])

            setShowPages(<div className="pages">
                {pagesSlice}
            </div>)
        }
    }, [allPage, currentPage])

    const handlerAddNode = (node) => {
        setNodes(prevState => prevState.length < 20 ? [node, ...prevState] : prevState)
    }

    const handlerEditNode = (node) => {
        setNodes(prevState => prevState.map(_node => _node.ID === node.ID ? node : _node))
    }

    const handlerDeleteNode = (nodeID) => {
        FetchRequest("DELETE", `/nodes/${nodeID}`)
            .then(response => {
                if (response.success && response.data) {
                    setNodes(prevState => prevState.filter(node => node.ID !== nodeID))
                }
            })
    }

    const handlerIndexNodes = () => {
        FetchRequest("GET", `/nodes/index`)
            .then(response => {
                console.log(response)
            })
    }

    return (
        <div className="contain nodes">
            {user.role.key !== "user" && canCreate && <>
                {modalCreate && <NodeModalCreate action={"create"} setState={setModalCreate} returnNode={handlerAddNode} defaultAddress={defaultAddress}/>}
                {modalEdit.State && <NodeModalCreate action={"edit"} setState={(state) => setModalEdit(prevState => ({...prevState, State: state}))} editNodeID={modalEdit.EditNodeID} returnNode={handlerEditNode}/>}
                <div className="buttons">
                    <button onClick={handlerIndexNodes}>
                        <FontAwesomeIcon icon={faRotate} /> Индексация
                    </button>
                    <button onClick={() => setModalCreate(true)}>
                        <FontAwesomeIcon icon={faPlus}/> Добавить узел
                    </button>
                </div>
            </>}
            {houseID === 0 && <input className="search" placeholder={"Поиск..."} type="text" value={search} onChange={handlerSearch}/>}
            {nodes.length > 0 ?
                <table>
                    <thead>
                    <tr className={"row-type-2"}>
                        <th>ID</th>
                        <th>Название</th>
                        {houseID > 0 ? "" : <th>Адрес</th>}
                        <th>Владелец</th>
                        <th>Район</th>
                        <th>Тип</th>
                        <th></th>
                    </tr>
                    </thead>
                    <tbody>
                    {isLoaded && nodes.map((node, index) => (
                        <tr key={index} className={index % 2 === 0 ? 'row-type-1' : 'row-type-2'}>
                            <td>{node.ID}</td>
                            <td>{node.Name}</td>
                            {houseID > 0 ? "" : <td>{`${node.Address.street.type.short_name} ${node.Address.street.name}, ${node.Address.house.type.short_name} ${node.Address.house.name}`}</td>}
                            <td>{node.Owner.Value}</td>
                            <td>{node.Zone.String}</td>
                            <td>{node.IsPassive ? <span className="bg-red">Пассивный</span> : <span className="bg-green">Активный</span>}</td>
                            {action === "select" ?
                                <td>
                                    <FontAwesomeIcon icon={selectNode?.ID === node.ID ? faCircleDot : faCircle} title="Выбрать" onClick={() => selectFunction(node)}/>
                                </td>
                            :
                                <td>
                                    <FontAwesomeIcon icon={faEye} className="eye" title="Просмотр" onClick={() => navigate(`/nodes/view/${node.ID}`)}/>
                                    {user.role.key !== "user" && <FontAwesomeIcon icon={faPen} title="Редактировать" onClick={() => setModalEdit({State: true, EditNodeID: node.ID})}/>}
                                    {user.role.key === "admin" && <FontAwesomeIcon icon={faTrash} className="delete" title="Удалить" onClick={() => handlerDeleteNode(node.ID)}/>}
                                </td>
                            }
                        </tr>
                    ))}
                    </tbody>
                </table>
                :
                <div className="empty">Нет узлов</div>
            }
            <div className="pagination">
                <div className="start" onClick={() => setCurrentPage(1)}><FontAwesomeIcon icon={faAnglesLeft}/>
                </div>
                <div className="back" onClick={() => {
                    if (currentPage - 1 > 0) {
                        setCurrentPage(prevState => prevState - 1)
                    }
                }}><FontAwesomeIcon icon={faAngleLeft}/></div>

                {showPages}

                <div className="next" onClick={() => {
                    if (currentPage + 1 <= allPage.length) {
                        setCurrentPage(prevState => prevState + 1)
                    }
                }}><FontAwesomeIcon icon={faAngleRight}/></div>
                <div className="end" onClick={() => setCurrentPage(allPage.length)}><FontAwesomeIcon
                    icon={faAnglesRight}/></div>
            </div>
        </div>
    )
}

export default NodesTable