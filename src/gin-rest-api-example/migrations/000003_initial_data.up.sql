-- admin password is "admin1"
-- user1 password is "user1"

INSERT INTO accounts (username, email, password, created_at, updated_at) VALUES
('admin', 'admin@email.com', '$2a$10$qUq8XJ.RcxcQvaEipicm/OLVPp5AjJjoPigj4vlRU579Xz0SkZwqu', now(), now()),
('user1', 'user1@email.com', '$2a$10$lsYsLv8nGPM0.R.ft4sgpe3OP7..KL3ZJqqhSVCKTEnSCMUztoUcW', now(), now());


INSERT INTO GameMode (id, mode) VALUES
  (1, 'Battle Royal'),
  (2, 'Deathmatch'),
  (3, 'Capture the Flag'),
  (4, 'Competitive'),
  (5, 'Unranked');


INSERT INTO GameModePopularity (areacode, gamemodeid, popularity) VALUES
  ('123', 1, 80),
  ('123', 2, 70),
  ('456', 3, 90),
  ('456', 4, 85),
  ('789', 5, 75),
  ('789', 4, 97),
  ('789', 3, 28),
  ('123', 5, 100),
  ('123', 4, 99);




