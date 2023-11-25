CREATE OR REPLACE FUNCTION regexp_match_count(search_text TEXT, pattern TEXT)
    RETURNS INTEGER AS $$
BEGIN
    RETURN (SELECT COUNT(*) FROM regexp_matches(search_text, pattern, 'g'));
END;
$$ LANGUAGE plpgsql;