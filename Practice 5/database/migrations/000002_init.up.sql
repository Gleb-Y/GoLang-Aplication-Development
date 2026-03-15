CREATE TABLE IF NOT EXISTS users (
    id         SERIAL       PRIMARY KEY,
    name       VARCHAR(255) NOT NULL,
    email      VARCHAR(255) NOT NULL DEFAULT '',
    gender     VARCHAR(10)  NOT NULL DEFAULT '',
    birth_date DATE         NOT NULL DEFAULT '2000-01-01'
);

CREATE TABLE IF NOT EXISTS user_friends (
    user_id   INTEGER REFERENCES users(id) ON DELETE CASCADE,
    friend_id INTEGER REFERENCES users(id) ON DELETE CASCADE,
    PRIMARY KEY (user_id, friend_id),
    CONSTRAINT no_self_friendship CHECK (user_id <> friend_id)
);

INSERT INTO users (name, email, gender, birth_date) VALUES
    ('Alice Johnson',   'alice@example.com',   'female', '1998-03-15'),
    ('Bob Smith',       'bob@example.com',     'male',   '1995-07-22'),
    ('Charlie Brown',   'charlie@example.com', 'male',   '1997-11-30'),
    ('Diana Prince',    'diana@example.com',   'female', '1999-05-10'),
    ('Edward Stark',    'edward@example.com',  'male',   '1996-08-04'),
    ('Fiona Green',     'fiona@example.com',   'female', '2000-01-20'),
    ('George White',    'george@example.com',  'male',   '1994-12-12'),
    ('Hannah Lee',      'hannah@example.com',  'female', '2001-06-25'),
    ('Ivan Torres',     'ivan@example.com',    'male',   '1993-09-09'),
    ('Julia Adams',     'julia@example.com',   'female', '2002-02-14'),
    ('Kevin Hart',      'kevin@example.com',   'male',   '1990-04-01'),
    ('Laura Palmer',    'laura@example.com',   'female', '1992-07-18'),
    ('Mike Ross',       'mike@example.com',    'male',   '1997-03-28'),
    ('Nancy Drew',      'nancy@example.com',   'female', '1998-10-05'),
    ('Oscar Wilde',     'oscar@example.com',   'male',   '1991-11-11'),
    ('Paula Abdul',     'paula@example.com',   'female', '1989-06-19'),
    ('Quinn Fisher',    'quinn@example.com',   'male',   '2003-08-30'),
    ('Rachel Green',    'rachel@example.com',  'female', '1996-04-22'),
    ('Sam Wilson',      'sam@example.com',     'male',   '1994-01-07'),
    ('Tina Turner',     'tina@example.com',    'female', '1999-12-31');

-- Alice (id=1) and Bob (id=2) share common friends: Charlie(3), Diana(4), Edward(5), Fiona(6)
INSERT INTO user_friends (user_id, friend_id) VALUES
    (1, 3), (3, 1),
    (1, 4), (4, 1),
    (1, 5), (5, 1),
    (1, 6), (6, 1),
    (2, 3), (3, 2),
    (2, 4), (4, 2),
    (2, 5), (5, 2),
    (2, 6), (6, 2),
    (1, 7), (7, 1),
    (2, 8), (8, 2),
    (9, 10), (10, 9),
    (11, 12), (12, 11),
    (13, 14), (14, 13),
    (15, 16), (16, 15),
    (17, 18), (18, 17),
    (19, 20), (20, 19);
