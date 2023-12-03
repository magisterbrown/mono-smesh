
from sqlalchemy import create_engine, text
from sqlalchemy.exc import IntegrityError
db_url = 'postgresql://magisterbrownie:post@localhost/auth_ranking'
engine = create_engine(db_url)
with engine.connect() as conn:
    try:
        conn.execute(text("INSERT INTO players (user_name, password) values ('art3', 'tar') RETURNING id"))
    except IntegrityError:
        import ipdb; ipdb.set_trace()
        print("dd")
    print(res.fetchall())
    res = conn.execute(text("SELECT * FROM players"))
    print(res.fetchall())
    conn.commit() 

