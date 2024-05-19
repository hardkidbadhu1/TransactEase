CREATE TABLE transactions (
    Transaction_ID INT PRIMARY KEY,
    Account_ID INT,
    OperationType_ID INT,
    Amount DECIMAL(10, 2),
    EventDate DATE DEFAULT NOW(),
    FOREIGN KEY (Account_ID) REFERENCES accounts(Account_ID),
    FOREIGN KEY (OperationType_ID) REFERENCES operations_types(OperationType_ID)
);
