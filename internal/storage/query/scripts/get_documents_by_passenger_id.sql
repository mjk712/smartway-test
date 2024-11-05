SELECT
    document_id,document_type,document_number, passenger_id AS "document.passenger_id"
FROM document
WHERE passenger_id = $1;