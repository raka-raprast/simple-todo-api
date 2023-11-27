import psycopg2
from psycopg2.extras import RealDictCursor

db_params = {
    'host': 'host',
    'port': 'port',
    'user': 'user',
    'password': 'password',
    'dbname': 'dbname'
}

conn = psycopg2.connect(**db_params)
conn.autocommit = True

cursor = conn.cursor(cursor_factory=RealDictCursor)
