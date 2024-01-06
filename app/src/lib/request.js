import { backend_url } from './config.js';

export async function authenticatedFetch(endpoint, options = {}) {
    const token = localStorage.getItem("Authorization");
    if(token == null) window.location.href='/signup';
    if(!options.headers) options.headers = {}
    options.headers['Authorization'] = token
    let resp = await fetch(backend_url+endpoint, options);
    console.log(resp);
    //if (resp.status === 404) {
    //   window.location.href='/signup'; 
    //}
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
