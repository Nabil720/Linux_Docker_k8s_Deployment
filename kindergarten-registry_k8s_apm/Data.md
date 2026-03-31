```bash

🔍 Final status check:
  - Elasticsearch: ✅ Running
  - Kibana: ✅ Running
  - APM Server: ✅ Running (auth required)

🎉 Setup complete!

📋 Service URLs:
  - Elasticsearch: http://localhost:9200
  - Kibana: http://localhost:5601
  - APM Server: http://localhost:8200

🔐 Credentials:
  - Username: elastic
  - Password: HEtL6W7qxEUJcs20

🔧 APM Configuration:
  - APM Server URL: http://localhost:8200
  - Secret Token: B7n5dCdEDTDppEbm

✨ You can now configure your applications to send APM data to http://localhost:8200
```



```
# ডিফল্ট credentials দিয়ে MongoDB তে প্রবেশ করুন (কোনো auth ছাড়া)
kubectl exec -it mongo-6dbfcb95dc-7xtn6 -- mongosh


// admin database এ switch করুন
use admin;

// বর্তমান users দেখুন
show users;

// পুরনো user ডিলিট করুন (যদি থাকে)
try {
  db.dropUser("myUser");
} catch(e) {
  print("User may not exist");
}

// নতুন user তৈরি করুন সঠিক credentials দিয়ে
db.createUser({
  user: "myUser",
  pwd: "myPassword",
  roles: [
    { role: "readWrite", db: "kindergarten" },
    { role: "readWrite", db: "admin" },
    { role: "dbAdmin", db: "kindergarten" },
    { role: "userAdminAnyDatabase", db: "admin" }
  ]
});

// user তৈরি হয়েছে কিনা চেক করুন
show users;

// authentication test
db.auth("myUser", "myPassword");

// kindergarten database তৈরি করুন
use kindergarten;

// collections তৈরি করুন
db.createCollection("employees");
db.createCollection("students");
db.createCollection("teachers");

// collections চেক করুন
show collections;

// exit
exit


```