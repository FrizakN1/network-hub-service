import React, {useContext, useEffect, useState} from "react";
import ImageOpen from "./ImageOpen";
import UploadFile from "./UploadFile";
import {useParams} from "react-router-dom";
import FetchRequest from "../fetchRequest";
import * as mime from 'react-native-mime-types';
import {faEllipsisVertical, faFolderPlus, faTrash} from "@fortawesome/free-solid-svg-icons";
import {FontAwesomeIcon} from "@fortawesome/react-fontawesome";
import AuthContext from "../context/AuthContext";

const ImageTable = ({type}) => {
    const [activeTab, setActiveTab] = useState(1)
    const [images, setImages] = useState([1])
    const [archiveImages, setArchiveImages] = useState([])
    const [openImage, setOpenImage] = useState({
        State: false,
        ImagesType: "active",
        Index: 0,
    })
    const { id } = useParams()
    const [panelIndex, setPanelIndex] = useState(null)
    const { user } = useContext(AuthContext)

    const getImageSrc = (image) => {
        const fileMimeType = mime.lookup(image.Name)
        const decodedData = atob(image.Data);

        const byteCharacters = new Uint8Array(decodedData.length);
        for (let i = 0; i < decodedData.length; i++) {
            byteCharacters[i] = decodedData.charCodeAt(i);
        }

        return URL.createObjectURL(new Blob([byteCharacters], { type: fileMimeType }))
    }

    useEffect(() => {
        setImages([])

        FetchRequest("GET", `/${type}/${id}/images`, null)
            .then(response => {
                if (response.success && response.data != null) {
                    let _archiveImages = []
                    let _images = []

                    for (let image of response.data) {
                        image.Src = getImageSrc(image)
                        image.InArchive ? _archiveImages.push(image) : _images.push(image)
                    }

                    setImages(_images)
                    setArchiveImages(_archiveImages)
                }
            })
    }, [type, id]);

    const handlerAddImage = (image) => {
        image.Src = getImageSrc(image)
        setImages(prevState => [image, ...prevState])
    }

    const handlerArchiveImage = (file) => {
        setPanelIndex(null)

        FetchRequest("POST", "/files/archive", file)
            .then(response => {
                if (response.success && response.data != null) {
                    response.data.Src = getImageSrc(response.data)

                    if (response.data.InArchive) {
                        setImages(prevState => prevState.filter(file =>
                            file.ID !== response.data.ID
                        ))

                        setArchiveImages(prevState => [response.data, ...prevState].sort((a, b) => b.ID - a.ID))
                    } else {
                        setArchiveImages(prevState => prevState.filter(file =>
                            file.ID !== response.data.ID
                        ))

                        setImages(prevState => [response.data, ...prevState].sort((a, b) => b.ID - a.ID))
                    }
                }
            })
    }

    const handlerDeleteImage = (file) => {
        setPanelIndex(null)

        FetchRequest("POST", "/files/delete", file)
            .then(response => {
                if (response.success && response.data != null) {
                    if (response.data.InArchive) {
                        setArchiveImages(prevState => prevState.filter(file => file.ID !== response.data.ID))
                    } else {
                        setImages(prevState => prevState.filter(file => file.ID !== response.data.ID))
                    }
                }
            })
    }

    return (
        <div style={{paddingBottom: "20px"}}>
            {openImage.State && <ImageOpen setState={(state) => setOpenImage(prevState => ({...prevState, State: state}))} images={openImage.ImagesType === "active" ? images : archiveImages} currentIndex={openImage.Index}/>}
            {user.role.key !== "user" && <UploadFile returnFile={handlerAddImage} type={type} onlyImage={true}/>}
            <div className="contain tables">
                <div className="tabs">
                    <div className={activeTab === 1 ? "tab active" : "tab"} onClick={() => {setActiveTab(1); setPanelIndex(null)}}>Актуальные изображения</div>
                    <div className={activeTab === 2 ? "tab active" : "tab"} onClick={() => {setActiveTab(2); setPanelIndex(null)}}>Архивированные изображения</div>
                </div>
                {activeTab === 1 ?
                    images.length > 0 ?
                        <div className="images">
                            {images.map((image, index) => (
                                <div key={"image"+index} className="image">
                                    {user.role.key !== "user" && <FontAwesomeIcon icon={faEllipsisVertical} style={panelIndex === index && {color: "#ffffff"}} className="menu" onClick={() => setPanelIndex(prevState => prevState === index ? null : index)}/>}
                                    {user.role.key !== "user" && panelIndex === index ?
                                        <div className="menu-block">
                                            {user.role.key !== "user" && <div onClick={() => handlerArchiveImage(image)}><FontAwesomeIcon icon={faFolderPlus} title="Переместить в архив"/> Переметить в архив</div>}
                                            {user.role.key === "admin" && <div onClick={() => handlerDeleteImage(image)}><FontAwesomeIcon icon={faTrash} title="Удалить" /> Удалить</div>}
                                        </div>
                                        :
                                        <img src={image.Src} alt="" onClick={() => setOpenImage({State: true, Index: index, ImagesType: "active"})}/>
                                    }
                                </div>
                            ))}
                        </div>
                        :
                        <div className="empty">Нет изображений</div>
                    :
                    archiveImages.length > 0 ?
                        <div className="images">
                            {archiveImages.map((image, index) => (
                                <div key={"image"+index} className="image">
                                    {user.role.key !== "user" && <FontAwesomeIcon icon={faEllipsisVertical} style={panelIndex === index && {color: "#ffffff"}} className="menu" onClick={() => setPanelIndex(prevState => prevState === index ? null : index)}/>}
                                    {panelIndex === index ?
                                        <div className="menu-block">
                                            {user.role.key !== "user" && <div onClick={() => handlerArchiveImage(image)}><FontAwesomeIcon icon={faFolderPlus} title="Восстановить"/> Восстановить</div>}
                                            {user.role.key === "admin" && <div onClick={() => handlerDeleteImage(image)}><FontAwesomeIcon icon={faTrash} title="Удалить" /> Удалить</div>}
                                        </div>
                                        :
                                        <img src={image.Src} alt="" onClick={() => setOpenImage({State: true, Index: index, ImagesType: "archive"})}/>
                                    }
                                </div>
                            ))}
                        </div>
                        :
                        <div className="empty archive">Нет изображений</div>
                }
            </div>
        </div>
    )
}

export default ImageTable