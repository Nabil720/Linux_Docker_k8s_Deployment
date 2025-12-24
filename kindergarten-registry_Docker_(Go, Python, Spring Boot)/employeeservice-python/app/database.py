from pymongo import MongoClient
import os
import time

class MongoDB:
    _client = None
    _db = None
    
    @classmethod
    def connect(cls):
        if cls._client is None:
            connection_string = os.getenv('MONGODB_URI', 'mongodb://mongo:27017')
            database_name = os.getenv('DATABASE_NAME', 'kindergarten')
            
            print(f"Connecting to MongoDB: {connection_string}")
            print(f"Database: {database_name}")
            
            # Retry connection with exponential backoff
            max_retries = 5
            for attempt in range(max_retries):
                try:
                    cls._client = MongoClient(connection_string, serverSelectionTimeoutMS=5000)
                    cls._db = cls._client[database_name]
                    
                    # Test connection
                    cls._client.server_info()
                    print("Connected to MongoDB successfully!")
                    break
                except Exception as e:
                    if attempt < max_retries - 1:
                        wait_time = 2 ** attempt
                        print(f"Connection attempt {attempt + 1} failed: {str(e)}. Retrying in {wait_time} seconds...")
                        time.sleep(wait_time)
                    else:
                        print(f"Failed to connect to MongoDB after {max_retries} attempts: {str(e)}")
                        raise
        
        return cls._db
    
    @classmethod
    def get_collection(cls, collection_name):
        if cls._db is None:
            cls.connect()
        return cls._db[collection_name]