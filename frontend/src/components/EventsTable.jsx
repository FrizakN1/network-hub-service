import React, {useEffect, useState} from "react";
import FetchRequest from "../fetchRequest";
import {FontAwesomeIcon} from "@fortawesome/react-fontawesome";
import {
    faAngleLeft,
    faAngleRight,
    faAnglesLeft, faAnglesRight,
} from "@fortawesome/free-solid-svg-icons";
import {useParams} from "react-router-dom";

const EventsTable = ({from = ""}) => {
    const [events, setEvents] = useState([])
    const [isLoaded, setIsLoaded] = useState(false)
    const [currentPage, setCurrentPage] = useState(1)
    const [allPage, setAllPage] = useState([])
    const [showPages, setShowPages] = useState(null)
    const [count, setCount] = useState(0)
    const [activeTab, setActiveTab] = useState("only")
    const { id } = useParams()

    useEffect(() => {
        handlerGetEvents()
    }, [activeTab]);

    const handlerGetEvents = () => {
        let uri = `/events`
        let params = new URLSearchParams({
            offset: String((currentPage-1)*20),
        })
        
        if (from !== "") {
            uri = `/${from}/${id}/events/${activeTab}`
        }

        FetchRequest("GET", `${uri}?${params.toString()}`, null)
            .then(response => {
                console.log(response)
                if (response.success) {
                    setEvents(response.data.Events != null ? response.data.Events : [])
                    setCount(response.data.Count)
                    setIsLoaded(true)
                }
            })
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
        handlerGetEvents()

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

    return (
        <div className="contain tables">
            {from !== "" &&
                <div className="tabs">
                    <div className={activeTab === "only" ? "tab active" : "tab"} onClick={() => setActiveTab("only")}>
                        События только {from === "houses" ? "дома" : from === "nodes" ? "узла" : "оборудования"}
                    </div>
                    <div className={activeTab === "all" ? "tab active" : "tab"} onClick={() => setActiveTab("all")}>
                        События связанные с {from === "houses" ? "домом" : from === "nodes" ? "узлом" : "оборудованием"}
                    </div>
                </div>
            }
            {events.length > 0 ?
                <table className={"events"}>
                    <thead>
                    <tr className={"row-type-2"}>
                        <th>ID</th>
                        <th className="col2">Описание</th>
                        <th>Адрес</th>
                        <th>Узел</th>
                        <th>Оборудование</th>
                        <th>Пользователь</th>
                        <th>Дата</th>
                        <th></th>
                    </tr>
                    </thead>
                    <tbody>
                    {isLoaded && events.map((event, index) => (
                        <tr key={index} className={index % 2 === 0 ? 'row-type-1' : 'row-type-2'}>
                            <td>{event.ID}</td>
                            <td className="col2">{event.Description}</td>
                            <td>{`${event.Address.Street.Type.ShortName} ${event.Address.Street.Name}, ${event.Address.House.Type.ShortName} ${event.Address.House.Name}`}</td>
                            <td>{event.Node != null ? event.Node.Name : "-"}</td>
                            <td>{event.Hardware != null ? event.Hardware.Type.TranslateValue : "-"}</td>
                            <td>{event.User.name}</td>
                            <td>{new Date(event.CreatedAt * 1000).toLocaleString().slice(0, 17)}</td>
                            <td></td>
                        </tr>
                    ))}
                    </tbody>
                </table>
                :
                <div className="empty">Нет событий</div>
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

export default EventsTable