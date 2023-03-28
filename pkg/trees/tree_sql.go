package trees

const (
	treesList = `
	SELECT id, name, description, is_active, create_time, creator, external_id
	FROM tree_mobile
	LIMIT $1 OFFSET $2;`

	treesGet = `
	SELECT id, name, description, external_id, is_active, inactivation_time, inactivation_reason, comment, is_validated, id_validator,
			create_time, creator, last_modification_time, last_modification_user, ST_AsText(geom) as geom, tree_attributes
	FROM tree_mobile
	WHERE id = $1`
	
	treesGetMaxId = "SELECT MAX(id) FROM tree_mobile;"

	treesExist = "SELECT COUNT(*) FROM tree_mobile WHERE id = $1;" 

	treesCount = "SELECT COUNT(*) FROM tree_mobile;"

	treesCreate = `
	INSERT INTO tree_mobile
	(name, description, external_id, is_active, comment, create_time, creator, geom, tree_attributes) 
	VALUES ($1, $2, $3, $4, $5, CURRENT_TIMESTAMP, $6, ST_GeomFromText($7, 2056), $8)
	RETURNING id;`
	
	treesUpdate = `
	UPDATE tree_mobile
	SET name					= $1,
		description				= $2,
		external_id				= $3,
		is_active				= $4,
		inactivation_time		= $5,
		inactivation_reason		= $6,
		comment					= $7,
		is_validated			= $8,
		id_validator			= $9,
		last_modification_time 	= CURRENT_TIMESTAMP,
		last_modification_user	= $10,
		geom					= ST_GeomFromText($11, 2056),
		tree_attributes			= $12
	WHERE id = $13;`

	treesDelete = "DELETE FROM tree_mobile WHERE id = $1;"

	treesSearchByName = "SELECT * FROM tree_mobile WHERE name LIKE $1;"

	treesIsActive = "SELECT isactive FROM tree_mobile WHERE id = $1;"

	treesCreateTable = `
	CREATE TABLE IF NOT EXISTS ${TABLE_NAME}
	(
	  id                      serial            CONSTRAINT tree_mobile_pk   primary key,
	  name                    text  not null constraint name_min_length check (length(btrim(name)) > 2),
	  description             text           constraint description_min_length check (length(btrim(description)) > 2),
	  external_id             int,
	  is_active               boolean default true not null,
	  inactivation_time       timestamp,
	  inactivation_reason     text,
	  comment                 text,
	  is_validated            boolean default false,
	  id_validator            int,
	  create_time             timestamp default now() not null,
	  creator                 integer  not null,
	  last_modification_time  timestamp,
	  last_modification_user  integer,
	  geom                    geometry(Point,2056)  not null,
	  tree_attributes         jsonb not null
	);
	ALTER TABLE tree_mobile OWNER TO postgres;
	COMMENT ON TABLE tree_mobile is 'tree_mobile is the main table of the sanarbo application';`
)