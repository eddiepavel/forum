CREATE TABLE connection(
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    post_id INTEGER,
    user_liked VARCHAR(255),
    FOREIGN KEY(username) REFERENCES user(username),
    FOREIGN KEY(post_id) REFERENCES post(id)
)