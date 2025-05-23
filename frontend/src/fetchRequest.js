import API_DOMAIN from "./config";

const FetchRequest = async (method, url, body = null) => {
    try {
        let options = {
            method: method,
            headers: {
                "Authorization": `Bearer ${localStorage.getItem("token")}`,
                "Content-Type": "application/json",
            }
        }

        if (body != null) {
            options = {
                ...options,
                body: JSON.stringify(body)
            }
        }

        const response = await fetch(API_DOMAIN+url, options)

        if (response.status === 401) {
            localStorage.removeItem("token");
            window.location.href = "/login";
            return{ success: false }
        }

        if (response.status === 403) {
            window.location.href = "/";
            return { success: false }
        }

        const data = await response.json();

        if (response.status === 400 || response.status === 500) {
            console.error(data.error)
            return { success: false }
        }

        return { success: true, data };
    } catch (error) {
        console.error(error)
        return { success: false, error };
    }
}

export default FetchRequest