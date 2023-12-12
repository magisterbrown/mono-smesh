export async function authenticatedFetch(url, options = {}) {
    let resp = await fetch(url, {... options});
    if (resp.status === 404) {
       window.location.href='/signup'; 
    }
    return resp
}
