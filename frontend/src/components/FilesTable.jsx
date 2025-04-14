import React, {useEffect, useState} from "react";
import {FontAwesomeIcon} from "@fortawesome/react-fontawesome";
import {
    faBoxArchive,
    faDownload, faFile,
    faFileExcel, faFileImage,
    faFileLines,
    faFilePdf,
    faFileWord,
    faFileZipper, faFolderMinus, faFolderPlus, faTrash, faTrashArrowUp
} from "@fortawesome/free-solid-svg-icons";
import API_DOMAIN from "../config";
import {useParams} from "react-router-dom";
import * as mime from 'react-native-mime-types';
import UploadFile from "./UploadFile";
import fetchRequest from "../fetchRequest";
import FetchRequest from "../fetchRequest";

const FilesTable = () => {
    const [files, setFiles] = useState([])
    const [archiveFiles, setArchiveFiles] = useState([])
    const [activeTab, setActiveTab] = useState(1)
    const { houseID } = useParams()

    useEffect(() => {
        setFiles([])

        FetchRequest("GET", `/get_files/${houseID}`, null)
            .then(response => {
                if (response.success && response.data != null) {
                    let _archiveFiles = []
                    let _files = []

                    for (let file of response.data) {
                        file.InArchive ? _archiveFiles.push(file) : _files.push(file)
                    }

                    setFiles(_files)
                    setArchiveFiles(_archiveFiles)
                }
            })

        // fetch(`${API_DOMAIN}/get_files/${houseID}`)
        //     .then(response => response.json())
        //     .then(data => {
        //         if (data != null) {
        //             let _archiveFiles = []
        //             let _files = []
        //
        //             for (let file of data) {
        //                 file.InArchive ? _archiveFiles.push(file) : _files.push(file)
        //             }
        //
        //             setFiles(_files)
        //             setArchiveFiles(_archiveFiles)
        //         }
        //     })
        //     .catch(error => console.error(error))
    }, [houseID]);

    const handlerDownloadFile = (file) => {
        const decodedData = atob(file.Data);
        const fileType = mime.lookup(file.Name)

        const byteCharacters = new Uint8Array(decodedData.length);
        for (let i = 0; i < decodedData.length; i++) {
            byteCharacters[i] = decodedData.charCodeAt(i);
        }

        const blob = new Blob([byteCharacters], { type: fileType });
        const url = URL.createObjectURL(blob);
        const a = document.createElement('a');

        a.href = url;
        a.download = file.Name;
        a.click();

        URL.revokeObjectURL(url);
    }

    const handlerArchiveFile = (file) => {
        // let options = {
        //     method: "POST",
        //     body: JSON.stringify(file)
        // }

        FetchRequest("POST", "/archive_file", file)
            .then(response => {
                if (response.success && response.data != null) {
                    if (response.data.InArchive) {
                        setFiles(prevState => prevState.filter(file =>
                            file.ID !== response.data.ID
                        ))

                        let updatedFiles = [response.data, ...archiveFiles]

                        setArchiveFiles(updatedFiles.sort((a, b) => b.ID - a.ID))
                    } else {
                        setArchiveFiles(prevState => prevState.filter(file =>
                            file.ID !== response.data.ID
                        ))

                        let updatedFiles = [response.data, ...files]

                        setFiles(updatedFiles.sort((a, b) => b.ID - a.ID))
                    }
                }
            })

        // fetch(`${API_DOMAIN}/archive_file`, options)
        //     .then(response => response.json())
        //     .then(data => {
        //         if (data != null) {
        //            if (data.InArchive) {
        //                setFiles(prevState => prevState.filter(file =>
        //                    file.ID !== data.ID
        //                ))
        //
        //                let updatedFiles = [data, ...archiveFiles]
        //
        //                setArchiveFiles(updatedFiles.sort((a, b) => b.ID - a.ID))
        //            } else {
        //                setArchiveFiles(prevState => prevState.filter(file =>
        //                    file.ID !== data.ID
        //                ))
        //
        //                let updatedFiles = [data, ...files]
        //
        //                setFiles(updatedFiles.sort((a, b) => b.ID - a.ID))
        //            }
        //         }
        //     })
        //     .catch(error => console.error(error))
    }

    const handlerDeleteFile = (file) => {
        // let options = {
        //     method: "POST",
        //     body: JSON.stringify(file)
        // }

        FetchRequest("POST", "/delete_file", file)
            .then(response => {
                if (response.success && response.data != null) {
                    if (response.data.InArchive) {
                        setArchiveFiles(prevState => prevState.filter(file => file.ID !== response.data.ID))
                    } else {
                        setFiles(prevState => prevState.filter(file => file.ID !== file.ID))
                    }
                }
            })

        // fetch(`${API_DOMAIN}/delete_file`, options)
        //     .then(response => response.json())
        //     .then(data => {
        //         if (data != null) {
        //             if (data.InArchive) {
        //                 setArchiveFiles(prevState => prevState.filter(file => file.ID !== data.ID))
        //             } else {
        //                 setFiles(prevState => prevState.filter(file => file.ID !== file.ID))
        //             }
        //         }
        //     })
        //     .catch(error => console.error(error))
    }

    const getIcon = (fileName) => {
        let icon
        let fileMimeType = mime.lookup(fileName)

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

        return icon
    }

    return (
        <div>
            <UploadFile setFiles={setFiles}/>
            <div className="contain tables">
                <div className="tabs">
                    <div className={activeTab === 1 ? "tab active" : "tab"} onClick={() => setActiveTab(1)}>Актуальные файлы</div>
                    <div className={activeTab === 2 ? "tab active" : "tab"} onClick={() => setActiveTab(2)}>Архивированные файлы</div>
                </div>
                {activeTab === 1 ?
                    files.length > 0 ? (
                                <table>
                                    <thead>
                                    <tr className={"row-type-2"}>
                                        <th className={"col1"}>Название файла</th>
                                        <th className={"col2"}>Дата загрузки</th>
                                        <th className={"col3"}></th>
                                    </tr>
                                    </thead>
                                    <tbody>
                                    {files.map((file, index) => (
                                        <tr key={index} className={index % 2 === 0 ? 'row-type-1' : 'row-type-2'}>
                                            <td className={"col1"}>{getIcon(file.Name)} {file.Name}</td>
                                            <td className={"col2"}>{new Date(file.UploadAt * 1000).toLocaleString().slice(0, 17)}</td>
                                            <td className={"col3"}>
                                                <FontAwesomeIcon icon={faDownload} title="Скачать" onClick={() => handlerDownloadFile(file)}/>
                                                <FontAwesomeIcon icon={faFolderPlus} title="Переместить в архив" onClick={() => handlerArchiveFile(file)}/>
                                            </td>
                                        </tr>
                                    ))}
                                    </tbody>
                                </table>
                            )
                            :
                            <div className="empty">Нет файлов</div>
                    :
                    archiveFiles.length > 0 ? (
                            <table className="archive">
                                <thead>
                                <tr className={"row-type-2"}>
                                    <th className={"col1"}>Название файла</th>
                                    <th className={"col2"}>Дата загрузки</th>
                                    <th className={"col3"}></th>
                                </tr>
                                </thead>
                                <tbody>
                                {archiveFiles.map((file, index) => (
                                    <tr key={index} className={index % 2 === 0 ? 'row-type-1' : 'row-type-2'}>
                                        <td className={"col1"}>{getIcon(file.Name)} {file.Name}</td>
                                        <td className={"col2"}>{new Date(file.UploadAt * 1000).toLocaleString().slice(0, 17)}</td>
                                        <td className={"col3"}>
                                            <FontAwesomeIcon icon={faDownload} title="Скачать" onClick={() => handlerDownloadFile(file)}/>
                                            <FontAwesomeIcon icon={faFolderMinus} title="Востановить" onClick={() => handlerArchiveFile(file)}/>
                                            <FontAwesomeIcon icon={faTrash} title="Удалить" onClick={() => handlerDeleteFile(file)}/>
                                        </td>
                                    </tr>
                                ))}
                                </tbody>
                            </table>
                        )
                        :
                        <div className="empty archive">Нет файлов</div>
                }
            </div>
        </div>
    )
}

export default FilesTable