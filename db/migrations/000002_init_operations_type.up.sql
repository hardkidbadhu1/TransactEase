CREATE TABLE IF NOT EXISTS operations_types (
    operation_type_id INT PRIMARY KEY,
    description VARCHAR(255) UNIQUE
);
