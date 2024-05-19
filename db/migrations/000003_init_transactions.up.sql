CREATE TABLE transactions (
    transaction_id SERIAL PRIMARY KEY,
    account_id INT,
    operation_type_id INT,
    amount DECIMAL(10, 2),
    event_date TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (account_id) REFERENCES accounts(account_id),
    FOREIGN KEY (operation_type_id) REFERENCES operations_types(operation_type_id)
);
