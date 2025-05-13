CREATE OR REPLACE FUNCTION insert_user_workspace()
RETURNS TRIGGER AS $$
BEGIN
	raise notice 'Trigger fired: Adding user_id %, workspace_id %', NEW.creater_id, NEW.id;
	insert into user_workspaces (user_id, workspace_id, user_member)
	values (NEW.creater_id, NEW.id, 'Owner');
	return NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trigger_insert_user_workspace
AFTER INSERT ON workspaces
FOR EACH ROW
EXECUTE FUNCTION insert_user_workspace();
