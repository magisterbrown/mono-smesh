from typing import Annotated

from fastapi import FastAPI, Header, HTTPException, Response, Body
from sqlalchemy import create_engine, text
from sqlalchemy.exc import IntegrityError, NoResultFound

from sqlalchemy.orm import sessionmaker
from sqlalchemy import create_engine

from models import Base, User, Token

from pydantic import BaseModel
from sqlalchemy import select
from configs.conf import db_url
import uuid
import hashlib


from fastapi.middleware.cors import CORSMiddleware
app = FastAPI()

app.add_middleware(
    CORSMiddleware,
    allow_credentials=True,
    allow_origins=['*'],
    allow_methods=["GET", "POST"],
    allow_headers=['*']
)

class Credentials(BaseModel):
    username: str
    password: str

subdomain = "/auth"

engine = create_engine(db_url)
Session = sessionmaker(autocommit=False, autoflush=False, bind=engine)

Base.metadata.create_all(engine)

def hasher(inp: str) -> str:
    return hashlib.sha512(inp.encode()).hexdigest()

@app.post(f"{subdomain}/signup")
def signup(cred: Credentials):
    db = Session()
    db_user = User(user_name=cred.username, password=hasher(cred.password))
    try:
        db.add(db_user)

        secret = uuid.uuid4().hex
        token = Token(user=db_user, token=hasher(secret))
        db.add(token)
        db.commit()
        return {"Authorization": secret}
    except IntegrityError:
        raise HTTPException(status_code=409, detail="User already exists")
 

@app.post(f"{subdomain}/login")
def login(cred: Credentials):
    db = Session()
    try:
        user = db.execute(select(User).where(User.user_name==cred.username, User.password==hasher(cred.password))).one()[0]
    except NoResultFound:
        raise HTTPException(status_code=401, detail="Wrong user name or password")

    secret = uuid.uuid4().hex
    token = Token(user=user, token=hasher(secret))
    db.add(token)
    db.commit()

    return {"Authorization": secret}

# TODO: when header is missing return 401
@app.api_route(subdomain, methods=["GET"])
def authentify(authorization: Annotated[str , Header()]):
    db = Session()
    if authorization is None:
        raise HTTPException(status_code=401)
    try:
        token = db.execute(select(Token).where(Token.token==hasher(authorization))).one()[0]
    except NoResultFound:
        raise HTTPException(status_code=401)
    return Response(status_code=200, headers={"User-Name": token.user.user_name})

