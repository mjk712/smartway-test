INSERT INTO passenger(last_name, middle_name, first_name)
VALUES('Matukov', 'Gigorish', 'Valerich'),
      ('Ostapov', 'Bender', 'Sishish'),
      ('Hulkenberg', 'Lubava', 'Gergova'),
      ('Mirka', 'Vassya', 'Tarasovna'),
      ('Trapov', 'Felix', 'Ignatievic');

INSERT INTO document(passenger_id, document_type, document_number)
VALUES(1,'passport', '34 25 876527'),
      (1,'spils', '378876527'),
      (1,'passport zagran', '52 76 983674'),
      (2,'passport', '98 86 522327'),
      (2,'bobr pass', '64 25 676384'),
      (3,'passport', '87 76 276382'),
      (4,'passport','50 60 304050'),
      (5,'passport','53 20 459234'),
      (5,'bobrpass','42 75 666777');

INSERT INTO flight_ticket(departure_point, destination_point, order_number, service_provider, departure_date, arrival_date, passenger_id)
VALUES('goida-town', 'Penza', '6924903', 'Nuke', '2024-12-23 13:40:04', '2024-12-24 1:45:12',1),
('anapa', 'Penza', '124237694', 'Pobeda', '2024-11-23 13:40:04', '2024-11-24 1:45:12',1),
('Guanjou', 'Pekin', '777666777', 'China Fly', '2025-01-23 13:40:04', '2025-01-24 1:45:12',2),
('Moskow', 'Dubai', '98765467', 'Fly Emirates', '2024-09-23 13:40:04', '2024-09-24 1:45:12',3),
('Moskow', 'Kumus', '1234984', 'Fly Kavkaz', '2024-08-23 13:40:04', '2024-08-24 1:45:12',4),
('Moskow', 'Tver', '8091643', 'Fly For Free', '2024-10-23 13:40:04', '2024-10-24 1:45:12',5);