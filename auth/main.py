from typing import Annotated

from fastapi import FastAPI, Header, HTTPException, Response
from sqlalchemy import create_engine, text
from sqlalchemy.exc import IntegrityError

from sqlalchemy.orm import sessionmaker
from sqlalchemy import create_engine

from models import Base, User, Token

app = FastAPI()
subdomain = "/auth"

db_url = 'postgresql://magisterbrownie:post@localhost/auth_ranking'
engine = create_engine(db_url)
Session = sessionmaker(autocommit=False, autoflush=False, bind=engine)

Base.metadata.create_all(engine)

@app.post(f"{subdomain}/signup")
def signup(username: str, password: str):
    db = Session()
    db_user = User(user_name=username, password=str(hash(password)))
    try:
        db.add(db_user)

        secret = uuid.uuid4().hex
        token = Token(user=db_user, token=hash(secret))
        db.commit()
        return {"Authorization": secret}
    except IntegrityError:
        raise HTTPException(status_code=409, detail="User already exists")
 

@app.post(f"{subdomain}/login")
def login():
    return {"Hello": "World"}

@app.api_route(subdomain, methods=["GET"])
def authentify(authorization: Annotated[list[str] | None, Header()] = None):
    if authorization is None:
        raise HTTPException(status_code=401)
    return Response(status_code=200, headers={"User-Name": "real DT"})

