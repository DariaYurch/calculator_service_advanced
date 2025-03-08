CREATE TABLE expressions(
    id CHAR(16) PRIMARY KEY,
    expression TEXT NOT NULL,
    status ENUM('PENDING', 'IN_PROGRESS', 'COMPLETED') NOT NULL DEFAULT 'PENDING',
    result FLOAT NULL
);

CREATE TABLE tasks(
    id CHAR(16) PRIMARY KEY,
    expression_id CHAR(16) NOT NULL,
    arg1 FLOAT,
    arg2 FLOAT,
    operation ENUM('+', '-', '*', '/') NOT NULL,
    status ENUM('PENDING', 'IN_PROGRESS', 'COMPLETED') NOT NULL DEFAULT 'PENDING',
    result FLOAT,
    FOREIGN KEY (expression_id) REFERENCES expressions(id) ON DELETE CASCADE
);

