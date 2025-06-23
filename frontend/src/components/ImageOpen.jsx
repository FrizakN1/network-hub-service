import React, {useEffect, useState} from "react";
import {faArrowLeft, faArrowRight} from "@fortawesome/free-solid-svg-icons";
import {FontAwesomeIcon} from "@fortawesome/react-fontawesome";
import {faSquareCaretLeft} from "@fortawesome/free-regular-svg-icons";

const ImageOpen = ({setState, images, currentIndex}) => {
    const [image, setImage] = useState("")
    const [index, setIndex] = useState(currentIndex)

    const handlerModalCreateClose = (e) => {
        if (e.target.className === "modal-window image") {
            setState(false)
        }
    }

    const handlerPressArrow = (e) => {
        if (e.key === 'ArrowLeft') {
            handlerPreviousImage()
        }
        if (e.key === "ArrowRight") {
            handlerNextImage()
        }
    };

    useEffect(() => {
        document.addEventListener('keydown', handlerPressArrow);
        return () => {
            document.removeEventListener('keydown', handlerPressArrow);
        };
    }, []);

    useEffect(() => {
        setImage(images[index].Src)
    }, [index])

    const handlerNextImage = () => {
        setIndex(prevState => {
            if (prevState + 1 <= images.length-1) {
               return prevState + 1
            }

            return 0
        })
    }

    const handlerPreviousImage = () => {
        setIndex(prevState => {
            if (prevState - 1 >= 0) {
                return prevState - 1
            }

            return images.length-1
        })
    }

    return (
        <div className={"modal-window image"} onMouseDown={handlerModalCreateClose}>
            <div onClick={handlerPreviousImage}><FontAwesomeIcon icon={faArrowLeft} /></div>
            <img src={image} alt=""/>
            <div onClick={handlerNextImage}><FontAwesomeIcon icon={faArrowRight} /></div>
        </div>
    )
}

export default ImageOpen