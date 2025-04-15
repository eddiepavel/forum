CREATE TABLE connection(
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    post_id INTEGER,
    user_id VARCHAR(255),
    FOREIGN KEY(user_id) REFERENCES user(id),
    FOREIGN KEY(post_id) REFERENCES post(id)
)