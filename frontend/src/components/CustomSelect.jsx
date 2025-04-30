import React, {useState} from "react";

const CustomSelect = ({placeholder, value, values, setValue}) => {
    const [isOlenList, setIsOpenList] = useState(false)

    const handlerSelectValue = (selectedValue) => {
        setIsOpenList(false)
        setValue(selectedValue)
    }

    return (
        <div className="custom-select">
            <div className="select-input" onClick={() => setIsOpenList(prevState => !prevState)}>{value === "" ? placeholder : value}</div>
            {isOlenList &&
                <ul className="list">
                    {values.map((_value, index) => (
                        <li key={"role"+index} onClick={() => handlerSelectValue(_value)}>{_value.TranslateValue}</li>
                    ))}
                </ul>
            }
        </div>
    )
}

export default CustomSelect