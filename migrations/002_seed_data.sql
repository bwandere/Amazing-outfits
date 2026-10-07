-- Migration: 002_seed_data.sql
-- Seed products and admin user

-- Default admin user: admin@amazingoutfits.co.ke / admin123
-- bcrypt hash for 'admin123': $2a$10$wT8hTqLCEu2fC9k1Q6F5k.gU3QWvA3k5K1Q6F5k.gU3QWvA3k5K1Q
INSERT INTO users (email, password_hash, role)
VALUES ('admin@amazingoutfits.co.ke', '$2a$10$3zO7Xyv6uQ/J/k9cZ3MvA.9Yt4hEw2iW9lH2v9Z0mE2dO3r4t5y6u', 'admin')
ON CONFLICT (email) DO NOTHING;

-- Seed Products
INSERT INTO products (id, name, brand, type, category, price, old_price, rating, reviews, description, tags)
VALUES
('volt-runner', 'Volt Runner', 'Amazing', 'sneaker', 'Running', 8500, 9800, 4.8, 124, 'Engineered mesh upper, gradient foam midsole and a responsive ride for long city miles.', '["featured", "popular"]'::jsonb),
('court-classic-low', 'Court Classic Low', 'Amazing', 'sneaker', 'Lifestyle', 7200, NULL, 4.7, 210, 'Clean tumbled leather, perforated toe box and a timeless low-cut silhouette.', '["featured", "popular"]'::jsonb),
('night-high', 'Night High', 'Amazing', 'sneaker', 'Basketball', 9900, NULL, 4.6, 88, 'All-black high-top with padded collar support and a blaze-orange traction outsole.', '["new", "featured"]'::jsonb),
('retro-wave', 'Retro Wave', 'Amazing', 'sneaker', 'Lifestyle', 8900, NULL, 4.9, 156, 'Chunky sculpted sole meets suede overlays in royal blue and burnt orange.', '["new", "popular", "featured"]'::jsonb),
('volt-runner-pro', 'Volt Runner Pro', 'Amazing', 'sneaker', 'Running', 11500, NULL, 4.7, 64, 'Our fastest runner with a carbon-infused plate and extra-bouncy foam.', '["new"]'::jsonb),
('street-court', 'Street Court', 'Amazing', 'sneaker', 'Basketball', 6800, 7900, 4.4, 97, 'Everyday hoops style with durable rubber cupsole for the street.', '["popular"]'::jsonb),
('dune-trainer', 'Dune Trainer', 'Amazing', 'sneaker', 'Training', 7600, NULL, 4.5, 41, 'Stable platform and grippy outsole for gym days and weekend walks.', '["new"]'::jsonb),
('blackout-mid', 'Blackout Mid', 'Amazing', 'sneaker', 'Training', 8200, NULL, 4.3, 52, 'Mid-top stealth trainer with lockdown lacing and cushioned heel.', '["popular"]'::jsonb),
('harambee-home', 'Harambee Stars Home 26', 'Amazing', 'jersey', 'National Teams', 3500, NULL, 4.8, 73, 'National-team inspired home kit with breathable knit and gold trims.', '["featured", "new"]'::jsonb),
('royal-blues-home', 'Royal Blues Home', 'Amazing', 'jersey', 'Clubs', 3200, NULL, 4.6, 58, 'Club home kit in royal blue with contrast orange collar and cuffs.', '["featured", "popular"]'::jsonb),
('red-stripes-classic', 'Red Stripes Classic', 'Amazing', 'jersey', 'Clubs', 3000, NULL, 4.5, 44, 'Iconic red-and-white vertical stripes with a classic crew neck.', '["popular"]'::jsonb),
('lions-away', 'Lions Away Kit', 'Amazing', 'jersey', 'National Teams', 3400, NULL, 4.4, 29, 'Lightweight away jersey with moisture-wicking fabric for match day.', '["new"]'::jsonb)
ON CONFLICT (id) DO NOTHING;

-- Seed Images
INSERT INTO product_images (product_id, url, display_order) VALUES
('volt-runner', '/assets/images/sneaker-1.jpg', 0),
('volt-runner', '/assets/images/sneaker-4.jpg', 1),
('volt-runner', '/assets/images/sneaker-3.jpg', 2),

('court-classic-low', '/assets/images/sneaker-2.jpg', 0),
('court-classic-low', '/assets/images/sneaker-1.jpg', 1),
('court-classic-low', '/assets/images/sneaker-4.jpg', 2),

('night-high', '/assets/images/sneaker-3.jpg', 0),
('night-high', '/assets/images/sneaker-2.jpg', 1),
('night-high', '/assets/images/sneaker-1.jpg', 2),

('retro-wave', '/assets/images/sneaker-4.jpg', 0),
('retro-wave', '/assets/images/sneaker-1.jpg', 1),
('retro-wave', '/assets/images/sneaker-2.jpg', 2),

('volt-runner-pro', '/assets/images/sneaker-1.jpg', 0),
('volt-runner-pro', '/assets/images/sneaker-3.jpg', 1),
('volt-runner-pro', '/assets/images/sneaker-4.jpg', 2),

('street-court', '/assets/images/sneaker-2.jpg', 0),
('street-court', '/assets/images/sneaker-3.jpg', 1),
('street-court', '/assets/images/sneaker-1.jpg', 2),

('dune-trainer', '/assets/images/sneaker-4.jpg', 0),
('dune-trainer', '/assets/images/sneaker-2.jpg', 1),
('dune-trainer', '/assets/images/sneaker-3.jpg', 2),

('blackout-mid', '/assets/images/sneaker-3.jpg', 0),
('blackout-mid', '/assets/images/sneaker-4.jpg', 1),
('blackout-mid', '/assets/images/sneaker-2.jpg', 2),

('harambee-home', '/assets/images/jersey-3.jpg', 0),
('harambee-home', '/assets/images/jersey-1.jpg', 1),
('harambee-home', '/assets/images/jersey-2.jpg', 2),

('royal-blues-home', '/assets/images/jersey-1.jpg', 0),
('royal-blues-home', '/assets/images/jersey-2.jpg', 1),
('royal-blues-home', '/assets/images/jersey-3.jpg', 2),

('red-stripes-classic', '/assets/images/jersey-2.jpg', 0),
('red-stripes-classic', '/assets/images/jersey-1.jpg', 1),
('red-stripes-classic', '/assets/images/jersey-3.jpg', 2),

('lions-away', '/assets/images/jersey-1.jpg', 0),
('lions-away', '/assets/images/jersey-3.jpg', 1),
('lions-away', '/assets/images/jersey-2.jpg', 2)
ON CONFLICT DO NOTHING;

-- Seed Sizes for Sneakers
INSERT INTO product_sizes (product_id, size, stock)
SELECT p.id, s.size, 10
FROM products p
CROSS JOIN (VALUES ('38'), ('39'), ('40'), ('41'), ('42'), ('43'), ('44'), ('45')) AS s(size)
WHERE p.type = 'sneaker'
ON CONFLICT DO NOTHING;

-- Seed Sizes for Jerseys
INSERT INTO product_sizes (product_id, size, stock)
SELECT p.id, s.size, 10
FROM products p
CROSS JOIN (VALUES ('S'), ('M'), ('L'), ('XL'), ('XXL')) AS s(size)
WHERE p.type = 'jersey'
ON CONFLICT DO NOTHING;
