import React, {useEffect, useState} from "react";
import {useParams} from "react-router-dom";
import API_DOMAIN from "../config";
import SearchInput from "./SearchInput";
import {FontAwesomeIcon} from "@fortawesome/react-fontawesome";
import {faDownload, faFileWord} from "@fortawesome/free-solid-svg-icons";
import FilesTable from "./FilesTable";
import UploadFile from "./UploadFile";
import fetchRequest from "../fetchRequest";
import FetchRequest from "../fetchRequest";

const HousePage = () => {
    const { houseID } = useParams()
    const [address, setAddress] = useState({})
    const [isLoaded, setIsLoaded] = useState(false)

    useEffect(() => {
        FetchRequest("GET", `/get_house/${houseID}`, null)
            .then(response => {
                if (response.success && response.data != null) {
                    setAddress(response.data)
                    setIsLoaded(true)
                }
            })

        // fetch(`${API_DOMAIN}/get_house/${houseID}`, {method: "GET"})
        //     .then(response => response.json())
        //     .then(data => {
        //         if (data != null) {
        //             setAddress(data)
        //             setIsLoaded(true)
        //         }
        //     })
        //     .catch(error => console.error(error))
    }, []);

    return (
        <section className="house">
            {isLoaded &&
                <div>
                    <SearchInput defaultValue={`${address.Street.Type.ShortName} ${address.Street.Name}, ${address.House.Type.ShortName} ${address.House.Name}`}/>
                    <FilesTable/>
                </div>
            }
        </section>
    )
}

export default HousePage