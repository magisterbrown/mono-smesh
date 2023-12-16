const backend_url = "http://localhost:8060"

export async function authenticatedFetch(endpoint, options = {}) {
    let resp = await fetch(backend_url+endpoint, {... options});
    if (resp.status === 404) {
       window.location.href='/signup'; 
    }
    return resp
}

export async function getToken(endpoint, username, password) {
    let resp = await fetch(backend_url+"/auth"+endpoint, { method: "POST", headers: {"Content-Type": "application/json"}, body:JSON.stringify({"username":username, "password": password})}) 
    if(resp.status == 200) {
        let body = await resp.json();
        localStorage.setItem("Authorization", JSON.stringify(body["Authorization"]));
        window.location.href = "/leaderboard";
    }
}
