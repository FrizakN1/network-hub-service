import React, {useEffect, useRef, useState} from "react";
import switcher from "ai-switcher-translit";
import {useNavigate} from "react-router-dom";
import FetchRequest from "../fetchRequest";

const SearchInput = ({defaultValue = "", action, returnAddress}) => {
    const [inputValue, setInputValue] = useState(defaultValue)
    const [suggestions, setSuggestions] = useState([])
    const [activeSuggestionIndex, setActiveSuggestionIndex] = useState(-1);
    const debounceTimer = useRef(0)
    const navigate = useNavigate()

    useEffect(() => {
        setInputValue(defaultValue)
    }, [defaultValue]);

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
                let params = new URLSearchParams({
                    search: value
                })

                FetchRequest("GET", `/houses/search?${params.toString()}`, null)
                    .then(response => {
                        if (response.success) {
                            setSuggestions(response.data?.Addresses != null ? response.data.Addresses : [])
                        }
                    })
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

                if (action === "select") {
                    returnAddress(suggestion)
                } else {
                    navigate(`/house/${suggestion.house.id}`)
                }

                setInputValue(`${suggestion.street.type.short_name} ${suggestion.street.name}, ${suggestion.house.type.short_name} ${suggestion.house.name}`);
                setSuggestions([]);
                setActiveSuggestionIndex(-1)
            } else if (inputValue.length > 0) {
                if (action === "select") {
                    returnAddress(suggestions[0])
                } else {
                    navigate('/result', {
                        state: { query: switcher.getSwitch(inputValue, {
                                input: {
                                    ".": ".",
                                    ",": ",",
                                }
                            })
                        }
                    });
                }
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
        if (action === "select") {
            returnAddress(suggestion)
        } else {
            navigate(`/house/${suggestion.house.id}`)
        }

        setInputValue(`${suggestion.street.type.short_name} ${suggestion.street.name}, ${suggestion.house.type.short_name} ${suggestion.house.name}`);
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
                            {suggestion.street.type.short_name} {suggestion.street.name}, {suggestion.house.type.short_name} {suggestion.house.name}
                        </li>
                    ))}
                </ul>
            )}
        </div>
    )
}

export default SearchInput