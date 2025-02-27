BEGIN;

CREATE TABLE IF NOT EXISTS puzzle_generation_logs (
		id BIGSERIAL PRIMARY KEY,
		request jsonb NOT NULL,
		fulfilled bool,
		error_status text,
		created_at timestamptz NOT NULL DEFAULT NOW(),
		completed_at timestamptz
);

CREATE TABLE IF NOT EXISTS puzzles (
		id BIGSERIAL PRIMARY KEY,
		uuid text UNIQUE NOT NULL,
		game_id bigint NOT NULL,
		turn_number integer NOT NULL,
		answer jsonb NOT NULL,
		author_id integer,
		lexicon text,
		before_text text,
		after_text text,
		rating jsonb NOT NULL,
		generation_id bigint NOT NULL,
		bucket_index integer NOT NULL,
		created_at timestamptz NOT NULL DEFAULT NOW(),
		FOREIGN KEY (game_id) REFERENCES games (id),
		FOREIGN KEY (author_id) REFERENCES users (id),
		FOREIGN KEY (generation_id) REFERENCES puzzle_generation_logs (id)
);

CREATE TABLE IF NOT EXISTS puzzle_tag_titles (
		id BIGSERIAL PRIMARY KEY,
	   	tag_title text NOT NULL
);

CREATE TABLE IF NOT EXISTS puzzle_tags (
	   	puzzle_id bigint NOT NULL,
	   	tag_id bigint NOT NULL,
		UNIQUE(puzzle_id, tag_id),
	   	FOREIGN KEY (puzzle_id) REFERENCES puzzles(id),
	   	FOREIGN KEY (tag_id) REFERENCES puzzle_tag_titles(id)
);

CREATE TABLE IF NOT EXISTS puzzle_attempts (
		puzzle_id bigint NOT NULL,
		user_id bigint NOT NULL,
		correct bool,
		attempts integer NOT NULL,
		new_user_rating jsonb,
		new_puzzle_rating jsonb,
		created_at timestamptz NOT NULL DEFAULT NOW(),
		updated_at timestamptz NOT NULL DEFAULT NOW(),
		UNIQUE(puzzle_id, user_id),
		FOREIGN KEY (puzzle_id) REFERENCES puzzles (id),
		FOREIGN KEY (user_id) REFERENCES users (id)
);

CREATE TABLE IF NOT EXISTS puzzle_votes (
		puzzle_id bigint NOT NULL,
		user_id bigint NOT NULL,
		vote integer,
		created_at timestamptz NOT NULL DEFAULT NOW(),
		UNIQUE(puzzle_id, user_id),
		FOREIGN KEY (puzzle_id) REFERENCES puzzles (id),
		FOREIGN KEY (user_id) REFERENCES users (id)
);

CREATE INDEX ON puzzle_attempts (updated_at);

INSERT INTO puzzle_tag_titles (tag_title) VALUES ('EQUITY') ON CONFLICT DO NOTHING;
INSERT INTO puzzle_tag_titles (tag_title) VALUES ('BINGO') ON CONFLICT DO NOTHING;
INSERT INTO puzzle_tag_titles (tag_title) VALUES ('ONLY_BINGO') ON CONFLICT DO NOTHING;
INSERT INTO puzzle_tag_titles (tag_title) VALUES ('BLANK_BINGO') ON CONFLICT DO NOTHING;
INSERT INTO puzzle_tag_titles (tag_title) VALUES ('NON_BINGO') ON CONFLICT DO NOTHING;
INSERT INTO puzzle_tag_titles (tag_title) VALUES ('POWER_TILE') ON CONFLICT DO NOTHING;
INSERT INTO puzzle_tag_titles (tag_title) VALUES ('BINGO_NINE_OR_ABOVE') ON CONFLICT DO NOTHING;
INSERT INTO puzzle_tag_titles (tag_title) VALUES ('CEL_ONLY') ON CONFLICT DO NOTHING;

COMMIT;