CREATE TABLE comment (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    content LONGTEXT NOT NULL,
    author INTEGER NOT NULL,
    post_id INTEGER NOT NULL,
    time DATETIME NOT NULL,
    upvotes INTEGER DEFAULT 0,
    downvotes INTEGER DEFAULT 0,
    vote_count INTEGER DEFAULT 0,
    FOREIGN KEY(author) REFERENCES user(id),
    CONSTRAINT fk_post_comment FOREIGN KEY (post_id) REFERENCES post(id) ON DELETE CASCADE
);