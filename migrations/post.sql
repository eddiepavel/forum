CREATE TABLE post (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    title TEXT NOT NULL,
    categories TEXT NOT NULL,
    content LONGTEXT NOT NULL,
    author INTEGER NOT NULL, 
    time DATETIME NOT NULL,
    upvotes INTEGER DEFAULT 0,
    downvotes INTEGER DEFAULT 0,
    vote_count INTEGER DEFAULT 0,
    FOREIGN KEY(author) REFERENCES user(id)
);


