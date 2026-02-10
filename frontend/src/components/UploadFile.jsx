import React, {useRef, useState} from "react";
import {
    faFile,
    faFileExcel,
    faFileImage,
    faFileLines,
    faFilePdf,
    faFileWord,
    faFileZipper
} from "@fortawesome/free-solid-svg-icons";
import {FontAwesomeIcon} from "@fortawesome/react-fontawesome";
import * as mime from 'react-native-mime-types'
import {useParams} from "react-router-dom";
import FetchRequest from "../fetchRequest";
import API_DOMAIN from "../config";

const UploadFile = ({returnFile, type, onlyImage = false}) => {
    const [uploadFile, setUploadFile] = useState({
        File: null,
        Name: "",
        Icon: null
    })
    const { id } = useParams()
    const fileInputRef = useRef(null);

    const handlerUploadFile = () => {
        if (uploadFile.File != null) {
            const reader = new FileReader();
            reader.onload = () => {
                const formData = new FormData();

                formData.append("file", uploadFile.File);
                formData.append("id", id);
                formData.append("type", type);
                formData.append("onlyImage", onlyImage);

                fetch(`${API_DOMAIN}/files/upload`, {
                    method: "POST",
                    body: formData,
                    headers: {"Authorization": `Bearer ${localStorage.getItem("token")}`}
                })
                    .then(response => response.json())
                    .then(response => {
                        if (response != null) {
                            returnFile(response)
                            setUploadFile({
                                File: null,
                                Name: "",
                                Icon: null
                            })
                            fileInputRef.current.value = ""
                        }
                    })
                    .catch(error => console.error(error))
            };
            reader.readAsDataURL(uploadFile.File);
        }
    }

    const handlerFileChange = (event) => {
        const file = event.target.files[0];

        if (!file) {return}

        const fileMimeType = mime.lookup(file.name)
        const allowedMimeTypes = onlyImage ?
            [
                "image/png",
                "image/jpeg"
            ]
            :
            [
                "application/msword",
                "application/vnd.openxmlformats-officedocument.wordprocessingml.document",
                "application/pdf",
                "text/plain",
                "application/vnd.ms-excel",
                "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet",
                "application/x-rar-compressed",
                "application/zip",
                "application/x-7z-compressed",
                "application/x-tar",
                "application/vnd.rar",
                "image/png",
                "image/jpeg",
                "application/octet-stream",
            ]
        if (!allowedMimeTypes.includes(fileMimeType)) {
            alert("Недопустимый формат файла. Пожалуйста, выберите файл другого формата.");
        }

        if (file.size > 10485760) {
            alert("Размер файла слишком большой.");
        }

        const reader = new FileReader();

        reader.onload = (e) => {
            let icon

            switch (fileMimeType) {
                case "application/msword":
                case "application/vnd.openxmlformats-officedocument.wordprocessingml.document":
                    icon = <FontAwesomeIcon icon={faFileWord} />; break
                case "application/pdf": icon = <FontAwesomeIcon icon={faFilePdf} />; break
                case "text/plain": icon = <FontAwesomeIcon icon={faFileLines} />; break
                case "application/vnd.ms-excel":
                case "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet":
                    icon = <FontAwesomeIcon icon={faFileExcel} />; break
                case "application/x-rar-compressed":
                case "application/zip":
                case "application/x-7z-compressed":
                case "application/x-tar":
                case "application/vnd.rar":
                    icon = <FontAwesomeIcon icon={faFileZipper} />; break
                case "image/png":
                case "image/jpeg":
                    icon = <FontAwesomeIcon icon={faFileImage} />; break
                default:
                    icon = <FontAwesomeIcon icon={faFile} />
            }

            setUploadFile({
                File: file,
                Name: file.name,
                Icon: icon,
            })
        }

        reader.readAsDataURL(file)
    }

    return (
        <div className="contain">
            <div className="upload">
                <label>
                    <input type="file" ref={fileInputRef} accept={onlyImage ? ".png, .jpeg, .jpg" : ".doc, .docx, .pdf, .txt, .xls, .xlsx, .rar, .zip, .7z, .tar, .png, .jpeg, .jpg .cdr"} max="10485760"
                    onInput={handlerFileChange}/>
                    {uploadFile.File != null ?
                        <div>{uploadFile.Icon} <span>{uploadFile.Name}</span></div>
                        :
                        <div><FontAwesomeIcon icon={faFile}/> <span>{onlyImage ? "Выберите изображение..." : "Выберите файл..."}</span></div>
                    }
                </label>
                <button onClick={handlerUploadFile}>Загрузить</button>
            </div>
        </div>
    )
}

export default UploadFile