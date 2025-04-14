import React, {useEffect, useState} from "react";
import {useLocation} from "react-router-dom";
import SearchInput from "./SearchInput";
import AddressesTable from "./AddressesTable";
import FetchRequest from "../fetchRequest";

const ResultPage = () => {
    const location = useLocation()
    const { query } = location.state || "";

    const [addresses, setAddresses] = useState([])
    const [count, setCount] = useState(0)
    const [offset, setOffset] = useState(1)
    const [isLoaded, setIsLoaded] = useState(false)

    useEffect(() => {
        if (query.length > 0 && typeof query === "string") {
            // let options = {
            //     method: "POST",
            //     body: JSON.stringify({
            //         Text: query,
            //         Limit: 20,
            //         Offset: (offset-1)*20,
            //     })
            // }
            
            let body = {
                Text: query,
                Limit: 20,
                Offset: (offset-1)*20,
            }

            FetchRequest("POST", "/search", body)
                .then(response => {
                    if (response.success) {
                        if (response.data != null) {
                            setAddresses(response.data?.Addresses || [])
                            setCount(response.data?.Count || 0)
                        }

                        setIsLoaded(true)
                    }
                })

            // fetch(`${API_DOMAIN}/search`, options)
            //     .then(response => response.json())
            //     .then(data => {
            //         if (data != null) {
            //             setAddresses(data?.Addresses || [])
            //             setCount(data?.Count || 0)
            //         }
            //         setIsLoaded(true)
            //     })
            //     .catch(error => console.error(error))
        }
    }, [query, offset]);

    return (
        <section className="result">
            <SearchInput defaultValue={query}/>
            {isLoaded && <AddressesTable addresses={addresses} count={count} setOffset={setOffset}/>}
        </section>
    )
}

export default ResultPage