CREATE TABLE votes (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    user_id INTEGER NOT NULL,
    post_id INTEGER,
    comment_id INTEGER,
    vote_type TEXT NOT NULL CHECK(vote_type IN ('upvote', 'downvote')),
    FOREIGN KEY(user_id) REFERENCES user(id),
    CONSTRAINT fk_post_votes FOREIGN KEY(post_id) REFERENCES post(id) ON DELETE CASCADE,
    FOREIGN KEY(comment_id) REFERENCES comment(id)
);