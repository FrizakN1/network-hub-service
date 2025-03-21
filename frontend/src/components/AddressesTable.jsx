import React, {useEffect, useState} from "react";
import {FontAwesomeIcon} from "@fortawesome/react-fontawesome";
import {faAngleLeft, faAngleRight, faAnglesLeft, faAnglesRight} from "@fortawesome/free-solid-svg-icons";
import {useNavigate} from "react-router-dom";

const AddressesTable = ({addresses, count, setOffset, h = null}) => {
    const [currentPage, setCurrentPage] = useState(1)
    const [allPage, setAllPage] = useState([])
    const [showPages, setShowPages] = useState(null)
    const navigate = useNavigate()

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
    }, [addresses, count]);

    useEffect(() => {
        setOffset(currentPage)

        console.log(allPage)

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
        <div className="contain">
            <div className="addresses-table">
                <h2>{h != null ? h : "Найдено совпадений: "}{count}</h2>
                <div className="list">
                    {addresses.length > 0 && addresses.map((address, index) => (
                        <div key={`address-${index}`} className="address" onClick={() => navigate(`/house/${address.House.ID}`)}>
                            <h3>{address.Street.Type.ShortName} {address.Street.Name}, {address.House.Type.ShortName} {address.House.Name}</h3>
                            <span>(файлов загружено: {address.FileAmount})</span>
                        </div>
                    ))}
                </div>
                <div className="pagination">
                    <div className="start" onClick={() => setCurrentPage(1)}><FontAwesomeIcon icon={faAnglesLeft}/>
                    </div>
                    <div className="back" onClick={() => {
                        if (currentPage - 1 > 0) {
                            setCurrentPage(currentPage - 1)
                        }
                    }}><FontAwesomeIcon icon={faAngleLeft}/></div>

                    {showPages}

                    <div className="next" onClick={() => {
                        if (currentPage + 1 <= allPage.length) {
                            setCurrentPage(currentPage + 1)
                        }
                    }}><FontAwesomeIcon icon={faAngleRight}/></div>
                    <div className="end" onClick={() => setCurrentPage(allPage.length)}><FontAwesomeIcon
                        icon={faAnglesRight}/></div>
                </div>
            </div>
        </div>
    )
}

export default AddressesTable