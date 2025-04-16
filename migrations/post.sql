CREATE TABLE post (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    title TEXT NOT NULL UNIQUE,
    categories TEXT NOT NULL,
    content VARCHAR(255),
    author INTEGER NOT NULL, 
    time DATETIME NOT NULL,
    upvotes INTEGER DEFAULT 0,
    downvotes INTEGER DEFAULT 0,
    FOREIGN KEY(author) REFERENCES user(id)
);


