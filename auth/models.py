from sqlalchemy.ext.declarative import declarative_base
from sqlalchemy.orm import sessionmaker, relationship, mapped_column
from sqlalchemy import Integer, String, ForeignKey

Base = declarative_base()

class User(Base):
    __tablename__ = "players"
    
    id = mapped_column(Integer, primary_key=True)
    user_name = mapped_column(String, unique=True)
    password = mapped_column(String)
    tokens = relationship("Token", back_populates="user_id")

class Token(Base):
    __tablename__ = "session_token"

    id = mapped_column(Integer, primary_key=True)
    user_id = mapped_column(ForeignKey('players.id'))
    token = mapped_column(String, unique=True)
    user = relationship("User", back_populates="tokens")


