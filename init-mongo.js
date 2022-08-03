db.createUser(
    {
        user: "trader",
        pwd: "trader",
        roles: [
            {
                role: "readWrite",
                db: "trader"
            }
        ]
    }
)