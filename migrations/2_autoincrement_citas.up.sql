-- Autoincrement for citas id on create
CREATE SEQUENCE citas_id_seq;

CREATE OR REPLACE FUNCTION citas_id_seq()
  RETURNS TRIGGER AS $$
BEGIN
    IF NEW.id IS NULL THEN
        NEW.id = nextval('citas_id_seq');
    END IF;
    RETURN NEW;
    END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER citas_id_trigger
  BEFORE INSERT ON citas
  FOR EACH ROW
  EXECUTE FUNCTION citas_id_seq();
