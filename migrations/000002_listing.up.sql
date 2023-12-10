CREATE TABLE IF NOT EXISTS listings(
    id bigserial PRIMARY KEY,
    created_at timestamp(0) NOT NULL DEFAULT NOW(),
    title text NOT NULL,
    description text NOT NULL,
    category text NOT NULL,
    roomCount integer NOT NULL,
    bathroomCount integer NOT NULL,
    guestCount integer NOT NULL,
    location text NOT NULL,
    price integer NOT NULL,
    owner_id bigint NOT NULL REFERENCES users(id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS images(
    id bigserial PRIMARY KEY,
    listing_id bigint NOT NULL REFERENCES listings(id) ON DELETE CASCADE,
    url text NOT NULL
);

CREATE TABLE IF NOT EXISTS bookings(
    id bigserial PRIMARY KEY,
    created_at timestamp(0) NOT NULL DEFAULT NOW(),
    listing_id bigint NOT NULL REFERENCES listings(id) ON DELETE CASCADE,
    guest_id bigint NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    check_in date NOT NULL,
    check_out date NOT NULL,
    price integer NOT NULL,
    total integer NOT NULL
);