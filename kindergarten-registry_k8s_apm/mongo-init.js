// MongoDB init script executed by the official mongo docker image
// It runs only when the database is first initialized (when /data/db is empty).

// Use the same database name as in your .env (default: kindergarten)
const dbName = "kindergarten";
const db = db.getSiblingDB(dbName);

function ensureCollection(name) {
  if (!db.getCollectionNames().includes(name)) {
    db.createCollection(name);
    print(`Created collection: ${name}`);
  } else {
    print(`Collection already exists: ${name}`);
  }
}

ensureCollection("employees");
ensureCollection("students");
ensureCollection("teachers");
