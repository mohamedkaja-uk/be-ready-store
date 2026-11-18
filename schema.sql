CREATE TABLE IF NOT EXISTS bundles (
    id VARCHAR(50) PRIMARY KEY,
    title VARCHAR(255) NOT NULL,
    category VARCHAR(100) NOT NULL,
    tagline TEXT NOT NULL,
    bundle_price NUMERIC(10, 2) NOT NULL,
    stock INT NOT NULL DEFAULT 50
);

CREATE TABLE IF NOT EXISTS items (
    id VARCHAR(50) PRIMARY KEY,
    bundle_id VARCHAR(50) REFERENCES bundles(id) ON DELETE CASCADE,
    name VARCHAR(255) NOT NULL,
    description TEXT NOT NULL,
    individual_price NUMERIC(10, 2) NOT NULL,
    image_url TEXT NOT NULL
);

CREATE TABLE IF NOT EXISTS orders (
    id VARCHAR(50) PRIMARY KEY,
    customer_name VARCHAR(255) NOT NULL,
    customer_email VARCHAR(255) NOT NULL,
    total_amount NUMERIC(10, 2) NOT NULL,
    status VARCHAR(50) NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS order_items (
    id SERIAL PRIMARY KEY,
    order_id VARCHAR(50) REFERENCES orders(id) ON DELETE CASCADE,
    bundle_id VARCHAR(50) REFERENCES bundles(id),
    quantity INT NOT NULL
);

-- Seed Initial UK Community Data
INSERT INTO bundles (id, title, category, tagline, bundle_price, stock) 
VALUES 
('1', 'British Weather Emergency Kit', 'British', 'Equipped for sudden UK power cuts, heavy rain, and winter storms.', 49.99, 50),
('2', 'Desi Monsoon & Winter Comfort Bundle', 'Indian', 'Essential comfort spices, immunity teas, and power backup readiness for UK households.', 39.99, 50),
('3', 'Traditional Desi Hospitality & Feast Kit', 'Pakistani', 'Everything you need on hand for unexpected guests, family feasts, and celebrations.', 55.00, 50)
ON CONFLICT (id) DO NOTHING;

INSERT INTO items (id, bundle_id, name, description, individual_price, image_url)
VALUES 
('i1', '1', 'LED Hand-Crank Torch', 'Tesco Essentials', 19.99, 'https://images.unsplash.com/photo-1543269865-cbf427effbad?w=300'),
('i2', '1', 'Thermal Mylar Blankets', 'Morrisons Outdoors', 14.99, 'https://images.unsplash.com/photo-1584308666744-24d5c474f2ae?w=300'),
('i3', '1', 'First Aid Travel Kit', 'Boots / Lidl Health', 18.99, 'https://images.unsplash.com/photo-1603398938378-e54eab446dde?w=300'),
('i4', '2', 'Authentic Chai & Masala Spice Pack', 'Curated cardamom, ginger, & loose leaf tea', 12.99, 'https://images.unsplash.com/photo-1576092768241-dec231879fc3?w=300'),
('i5', '2', 'Ayurvedic Herbal Immunity Kit', 'Turmeric, ashwagandha, and honey essentials', 16.99, 'https://images.unsplash.com/photo-1540420773420-3366772f4999?w=300'),
('i6', '2', 'Heavy Duty Multi-Plug Extension Lead', 'UK 13A surge-protected household power sharing', 14.99, 'https://images.unsplash.com/photo-1558494949-ef010cbdcc31?w=300'),
('i7', '3', 'Premium Basmati Rice & Lentil Combo', 'Aromatic aged basmati and staple daal pack', 18.99, 'https://images.unsplash.com/photo-1586201375761-83865001e31c?w=300'),
('i8', '3', 'Specialty Biryani Spice Box & Ghee', 'Authentic layered spices and pure ghee', 15.99, 'https://images.unsplash.com/photo-1596797038530-2c107229654b?w=300'),
('i9', '3', 'Instant Dessert & Vermicelli Kit', 'Sheer khurma and kheer preparation essentials', 12.99, 'https://images.unsplash.com/photo-1512621776951-a57141f2eefd?w=300'),
('i10', '3', 'Large Stainless Steel Serving Tray Set', 'Traditional hospitality serving essentials', 19.99, 'https://images.unsplash.com/photo-1578749556568-bc2c40e68b61?w=300')
ON CONFLICT (id) DO NOTHING;
