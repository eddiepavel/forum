CREATE TABLE comment (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    content VARCHAR(255) NOT NULL,
    author INTEGER NOT NULL,
    post_id INTEGER NOT NULL,
    time DATETIME NOT NULL,
    upvotes INTEGER DEFAULT 0,
    downvotes INTEGER DEFAULT 0,
    FOREIGN KEY(author) REFERENCES user(id),
    FOREIGN KEY(post_id) REFERENCES post(id)
);