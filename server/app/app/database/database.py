import databases
import sqlalchemy
# import os

# url = os.environ.get('DATABASE_URL')

# print(url)
# Dev: postgresql://postgres:postgres@127.0.0.1:5432
# url = 'postgresql://postgres:postgres@0.0.0.0:5432'
import os

DATABASE_USER = os.environ.get('DATABASE_USER')
DATABASE_PASS = os.environ.get('DATABASE_PASS')
DATABASE_IP = os.environ.get('DATABASE_URI')
DATABASE_PORT = os.environ.get('DATABASE_PORT')

# DATABASE_URL = 'postgresql://postgres:postgres@0.0.0.0:5432' if DATABASE_URL is None else DATABASE_URL
# DEV database currently set up.
# DATABASE_URL = f"{DATABASE_URL}"
# postgresql://cwyf-admin:6fbedeb5-0118-4224-a17e-e6c0daef6467@cwyf.database.windows.net:5432
DATABASE_URL = f'postgresql://{DATABASE_USER}:{DATABASE_PASS}@{DATABASE_IP}:{DATABASE_PORT}'
database = databases.Database(DATABASE_URL)

async def startup():
    # Start database on app startup.
    try:
        await database.connect()
        print("INFO:     Successfully connected to database.")
    except ConnectionRefusedError:
        print("ERROR:    Could not connect to database.")
        raise 

async def shutdown():
    # Stop database on shutdown.
    print("INFO:     Disconnecting from database.")
    await database.disconnect()

def provide_connection() -> databases.Database:
        return database