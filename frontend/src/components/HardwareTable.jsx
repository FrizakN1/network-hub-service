import React, {useContext, useEffect, useRef, useState} from "react";
import {FontAwesomeIcon} from "@fortawesome/react-fontawesome";
import {
    faAngleLeft,
    faAngleRight,
    faAnglesLeft, faAnglesRight,
    faEye,
    faPen,
    faPlus,
    faTrash
} from "@fortawesome/free-solid-svg-icons";
import {useNavigate} from "react-router-dom";
import FetchRequest from "../fetchRequest";
import HardwareModalCreate from "./HardwareModalCreate";
import AuthContext from "../context/AuthContext";

const HardwareTable = ({id = 0, type = "", canCreate = false}) => {
    const [modalCreate, setModalCreate] = useState(false)
    const [modalEdit, setModalEdit] = useState({
        State: false,
        EditHardwareID: null
    })
    const searchDebounceTimer = useRef(0)
    const [hardware, setHardware] = useState([])
    const [isLoaded, setIsLoaded] = useState(false)
    const [currentPage, setCurrentPage] = useState(1)
    const [allPage, setAllPage] = useState([])
    const [showPages, setShowPages] = useState(null)
    const [search, setSearch] = useState("")
    const [count, setCount] = useState(0)
    const navigate = useNavigate()
    const { user } = useContext(AuthContext)

    useEffect(() => {
        handlerGetHardware()
    }, []);

    const handlerGetHardware = (value = "") => {
        let uri = "/hardware"
        let params = new URLSearchParams({
            offset: String((currentPage-1)*20),
        })

        if (value.length > 0) {
            params.set("search", value)
            uri = "/hardware/search"
        }

        if (type !== "") uri =  `/${type}/${id}/hardware`

        FetchRequest("GET", `${uri}?${params.toString()}`, null)
            .then(response => {
                if (response.success) {
                    setHardware(response.data.Hardware != null ? response.data.Hardware : [])
                    setCount(response.data.Count)
                    setIsLoaded(true)
                }
            })
    }

    const handlerSearch = (e) => {
        setSearch(e.target.value)

        clearTimeout(searchDebounceTimer.current)

        searchDebounceTimer.current = setTimeout(() =>  {
            handlerGetHardware(e.target.value)
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
        handlerGetHardware(search)

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

    const handlerAddHardware = (_hardware) => {
        setHardware(prevState => prevState.length < 20 ? [_hardware, ...prevState] : prevState)
    }

    const handlerEditHardware = (_hardware) => {
        setHardware(prevState => prevState.map(rec => rec.ID === _hardware.ID ? _hardware : rec))
    }

    const handlerDeleteHardware = (hardwareID) => {
        FetchRequest("DELETE", `/hardware/${hardwareID}`, null)
            .then(response => {
                if (response.success && response.data) {
                    setHardware(prevState => prevState.filter(hardware => hardware.ID !== hardwareID))
                }
            })
    }

    return (
        <div className="contain hardware">
            {user.role.key !== "user" && canCreate && <>
                {modalCreate && <HardwareModalCreate action={"create"} setState={setModalCreate} returnHardware={handlerAddHardware}/>}
                {modalEdit.State && <HardwareModalCreate action={"edit"} setState={(state) => setModalEdit(prevState => ({...prevState, State: state}))} editHardwareID={modalEdit.EditHardwareID} returnHardware={handlerEditHardware}/>}
                <div className="buttons">
                    <button onClick={() => setModalCreate(true)}>
                        <FontAwesomeIcon icon={faPlus}/> Добавить оборудование
                    </button>
                </div>
            </>}
            {id === 0 && <input className="search" placeholder={"Поиск..."} type="text" value={search} onChange={handlerSearch}/>}
            {hardware.length > 0 ?
                <table>
                    <thead>
                    <tr className={"row-type-2"}>
                        <th>ID</th>
                        <th>Тип оборудования</th>
                        {type !== "node" && <th>Узел</th>}
                        {type !== "house" && <th>Адрес</th>}
                        <th>Модель</th>
                        <th>IP адрес</th>
                        <th></th>
                    </tr>
                    </thead>
                    <tbody>
                    {isLoaded && hardware.map((_hardware, index) => (
                        <tr key={index} className={index % 2 === 0 ? 'row-type-1' : 'row-type-2'}>
                            <td>{_hardware.ID}</td>
                            <td>{_hardware.Type.Value}</td>
                            {type !== "node" && <td>{_hardware.Node.Name}</td>}
                            {type !== "house" && <td>{`${_hardware.Node.Address.Street.Type.ShortName} ${_hardware.Node.Address.Street.Name}, ${_hardware.Node.Address.House.Type.ShortName} ${_hardware.Node.Address.House.Name}`}</td>}
                            <td>{_hardware.Type.Key === "switch" ? _hardware.Switch.Name : "-"}</td>
                            <td>{_hardware.Type.Key === "switch" && _hardware.IpAddress.Valid ? _hardware.IpAddress.String : "-"}</td>
                            <td>
                                <FontAwesomeIcon icon={faEye} className="eye" title="Просмотр" onClick={() => navigate(`/hardware/view/${_hardware.ID}`)}/>
                                {user.role.key !== "user" &&<FontAwesomeIcon icon={faPen} title="Редактировать" onClick={() => setModalEdit({State: true, EditHardwareID: _hardware.ID})}/>}
                                {user.role.key === "admin" && <FontAwesomeIcon icon={faTrash} className="delete" title="Удалить" onClick={() => handlerDeleteHardware(_hardware.ID)}/>}
                            </td>
                        </tr>
                    ))}
                    </tbody>
                </table>
                :
                <div className="empty">Нет оборудования</div>
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

export default HardwareTable