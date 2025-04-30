import React from "react";
import {FontAwesomeIcon} from "@fortawesome/react-fontawesome";
import {faCircleXmark} from "@fortawesome/free-solid-svg-icons";

const InputErrorDescription = ({text}) => {
    return (
        <div className="input-error">
            <i className="error-triangle"></i>
            <FontAwesomeIcon icon={faCircleXmark}/>
            <p>{text}</p>
        </div>
    )
}

export default InputErrorDescription