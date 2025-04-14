import React, {useRef, useState} from "react";
import switcher from "ai-switcher-translit";
import {useNavigate} from "react-router-dom";
import FetchRequest from "../fetchRequest";

const SearchInput = ({defaultValue = ""}) => {
    const [inputValue, setInputValue] = useState(defaultValue)
    const [suggestions, setSuggestions] = useState([])
    const [activeSuggestionIndex, setActiveSuggestionIndex] = useState(-1);
    const debounceTimer = useRef(0)
    const navigate = useNavigate()

    const handlerChange = (e) => {
        setInputValue(e.target.value)

        let value = switcher.getSwitch(e.target.value, {
            input: {
                ".": ".",
                ",": ",",
            }
        })

        clearTimeout(debounceTimer.current)

        if (value.length > 0) {
            debounceTimer.current = setTimeout(() => {
                // let options = {
                //     method: "POST",
                //     body: JSON.stringify({
                //         Text: value,
                //     })
                // }

                FetchRequest("POST", "/search", {Text: value})
                    .then(response => {
                        if (response.success) {
                            if (response.data != null) {
                                setSuggestions(response.data?.Addresses || [])
                            } else {
                                setSuggestions([])
                            }
                        }
                    })

                // fetch(`${API_DOMAIN}/search`, options)
                //     .then(response => response.json())
                //     .then(data => {
                //         if (data != null) {
                //             setSuggestions(data?.Addresses || [])
                //         } else {
                //             setSuggestions([])
                //         }
                //     })
                //     .catch(error => console.error(error))
            }, 300)
        } else {
            setSuggestions([])
        }
    }

    const handleKeyDown = (e) => {
        if (e.key === "ArrowDown") {
            // Перемещаемся вниз по списку
            setActiveSuggestionIndex((prevIndex) =>
                prevIndex < suggestions.length - 1 ? prevIndex + 1 : prevIndex
            );
        } else if (e.key === "ArrowUp") {
            // Перемещаемся вверх по списку
            setActiveSuggestionIndex((prevIndex) =>
                prevIndex > 0 ? prevIndex - 1 : 0
            );
        } else if (e.key === "Enter") {
            // Выбираем текущую подсказку
            if (activeSuggestionIndex >= 0) {
                let suggestion = suggestions[activeSuggestionIndex]

                navigate(`/house/${suggestion.House.ID}`)

                setInputValue(`${suggestion.Street.Type.ShortName} ${suggestion.Street.Name}, ${suggestion.House.Type.ShortName} ${suggestion.House.Name}`);
                setSuggestions([]);
                setActiveSuggestionIndex(-1)
            } else if (inputValue.length > 0) {
                navigate('/result', {
                    state: { query: switcher.getSwitch(inputValue, {
                            input: {
                                ".": ".",
                                ",": ",",
                            }
                        })
                    }
                });
                clearTimeout(debounceTimer.current)
                setSuggestions([]);
            }
        } else if (e.key === "Escape") {
            // Закрываем подсказки при нажатии Escape
            setSuggestions([]);
            setActiveSuggestionIndex(-1)
        }
    };

    const handleClickSuggestion = (suggestion) => {
        navigate(`/house/${suggestion.House.ID}`)

        setInputValue(`${suggestion.Street.Type.ShortName} ${suggestion.Street.Name}, ${suggestion.House.Type.ShortName} ${suggestion.House.Name}`);
        setSuggestions([]);
    };

    return (
        <div className="search">
            <input type="text" value={inputValue} placeholder="Поиск адреса..." onChange={handlerChange} onKeyDown={handleKeyDown}/>
            {suggestions.length > 0 && (
                <ul className="suggestions">
                    {suggestions.map((suggestion, index) => (
                        <li key={index} onClick={() => handleClickSuggestion(suggestion)}
                            className={index === activeSuggestionIndex ? "active" : ""}>
                            {suggestion.Street.Type.ShortName} {suggestion.Street.Name}, {suggestion.House.Type.ShortName} {suggestion.House.Name}
                        </li>
                    ))}
                </ul>
            )}
        </div>
    )
}

export default SearchInput