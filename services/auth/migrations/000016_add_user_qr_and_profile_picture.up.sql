ALTER TABLE auth.users
    ADD COLUMN qr_code             TEXT,
    ADD COLUMN profile_picture_url TEXT;