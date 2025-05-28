CREATE TABLE notification (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    target_id INTEGER NOT NULL,
    actor_id INTEGER NOT NULL,
    post_id INTEGER NOT NULL,
    type TEXT NOT NULL CHECK ( type IN ('comment', 'downvote', 'upvote')),
    content TEXT,
    read TINYINT DEFAULT 0,
    time DATETIME NOT NULL,
    FOREIGN KEY(target_id) REFERENCES user(id),
    FOREIGN KEY(actor_id) REFERENCES user(id),
    CONSTRAINT fk_post_id FOREIGN KEY(post_id) REFERENCES post(id) ON DELETE CASCADE
)