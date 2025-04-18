import API_DOMAIN from "./config";

const FetchRequest = async (method, url, body) => {
    try {
        let options = {
            method: method,
            headers: {"Authorization": `Bearer ${localStorage.getItem("token")}`}
        }

        if (method === "POST") {
            options = {
                ...options,
                body: JSON.stringify(body)
            }
        }

        const response = await fetch(API_DOMAIN+url, options)

        if (response.status === 401) {
            localStorage.removeItem("token");
            window.location.href = "/login";
        }

        if (response.status === 403) {
            window.location.href = "/";
        }

        const data = await response.json();

        return { success: true, data };
    } catch (error) {
        console.error(error)
        return { success: false, error };
    }
}

export default FetchRequest