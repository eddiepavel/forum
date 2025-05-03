CREATE TABLE post (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    title TEXT NOT NULL,
    categories TEXT NOT NULL,
    content TEXT NOT NULL, -- Changed LONGTEXT to TEXT
    author INTEGER NOT NULL, 
    time DATETIME NOT NULL,
    upvotes INTEGER DEFAULT 0,
    downvotes INTEGER DEFAULT 0,
    vote_count INTEGER DEFAULT 0,
    image_path TEXT,
    FOREIGN KEY(author) REFERENCES user(id)
    
);