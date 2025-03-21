import React, {useRef, useState} from "react";
import API_DOMAIN from "../config";
import switcher from "ai-switcher-translit";
import SearchInput from "./SearchInput";

const MainPage = () => {
    return (
        <section className="home">
            <SearchInput />
        </section>
    )
}

export default MainPage