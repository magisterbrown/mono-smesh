from typing import Annotated

from fastapi import FastAPI, Header, HTTPException, Response

app = FastAPI()
subdomain = "/auth"


@app.post(f"{subdomain}/signup")
def read_root():
    return {"Hello": "World"}

@app.post(f"{subdomain}/login")
def read_root():
    return {"Hello": "World"}

@app.api_route(subdomain, methods=["GET"])
def read_root(authorization: Annotated[list[str] | None, Header()] = None):
    if authorization is None:
        raise HTTPException(status_code=401)
    return Response(status_code=200, headers={"User-Name": "real DT"})

