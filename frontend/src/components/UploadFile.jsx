import React, {useState} from "react";
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
import API_DOMAIN from "../config";
import {useParams} from "react-router-dom";
import FetchRequest from "../fetchRequest";

const UploadFile = ({setFiles}) => {
    const [uploadFile, setUploadFile] = useState({
        File: null,
        Name: "",
        Icon: null
    })
    const { houseID } = useParams()

    const handlerUploadFile = () => {
        if (uploadFile.File != null) {
            const reader = new FileReader();
            reader.onload = () => {
                const formData = new FormData();

                formData.append("file", uploadFile.File);
                formData.append("houseID", houseID);

                // let options = {
                //     method: "POST",
                //     body: formData
                // }

                FetchRequest("POST", "/upload_file", formData)
                    .then(response => {
                        if (response.success && response.data != null) {
                            setFiles(prevState => [response.data, ...prevState])
                            setUploadFile({
                                File: null,
                                Name: "",
                                Icon: null
                            })
                        }
                    })
                // fetch(`${API_DOMAIN}/upload_file`, options)
                //     .then(response => response.json())
                //     .then(data => {
                //         if (data != null) {
                //             setFiles(prevState => [data, ...prevState])
                //             setUploadFile({
                //                 File: null,
                //                 Name: "",
                //                 Icon: null
                //             })
                //         }
                //     })
                //     .catch(error => console.error(error))
            };
            reader.readAsDataURL(uploadFile.File);
        }
    }

    const handlerFileChange = (event) => {
        const file = event.target.files[0];

        if (file) {
            const fileMimeType = mime.lookup(file.name)
            const allowedMimeTypes = [
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
                "image/png",
                "image/jpeg"
            ];
            if (allowedMimeTypes.includes(fileMimeType)) {
                if (file.size <= 10485760) {
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
                } else {
                    alert("Размер файла слишком большой.");
                }
            } else {
                alert("Недопустимый формат файла. Пожалуйста, выберите файл другого формата.");
            }
        }
    }

    return (
        <div className="contain">
            <div className="upload">
                <label>
                    <input type="file" accept=".doc, .docx, .pdf, .txt, .xls, .xlsx, .rar, .zip, .7z, .tar, .png, .jpeg, .jpg" max="10485760"
                    onChange={handlerFileChange}/>
                    {uploadFile.File != null ?
                        <div>{uploadFile.Icon} <span>{uploadFile.Name}</span></div>
                        :
                        <div><FontAwesomeIcon icon={faFile}/> <span>Выберите файл...</span></div>
                    }
                </label>
                <button onClick={handlerUploadFile}>Загрузить</button>
            </div>
        </div>
    )
}

export default UploadFile