package trees

/*
	List(offset, limit int) ([]*TreeList, error)

	// Get returns the object with the specified objects ID.
	Get(id int32) (*Tree, error)

	// GetMaxId returns the maximum value of objects id existing in store.
	GetMaxId() (int32, error)

	// Exist returns true only if a objects with the specified id exists in store.
	Exist(id int32) bool

	// Count returns the total number of objects.
	Count() (int32, error)

	// Create saves a new objects in the storage.
	Create(object Tree) (*Tree, error)

	// Update updates the objects with given ID in the storage.
	Update(id int32, object Tree) (*Tree, error)

	// Delete removes the objects with given ID from the storage.
	Delete(id int32) error

	// SearchTreesByName list of existing objects where the name contains the given search pattern or err if not found
	SearchTreesByName(pattern string) ([]*TreeList, error)

	// IsTreeActive returns true if the object with the specified id has the is_active attribute set to true
	IsTreeActive(id int32) bool
*/


const (
	/*
	treesList = `
	SELECT 	thi.idthing AS id, thi.name, thi.description, thi.isactive AS is_active, thi.idcreator AS creator, thi.datecreated AS create_time, null AS external_id
	FROM thing thi
	WHERE thi.idtypething=74
	ORDER BY thi.idthing
	LIMIT $1 OFFSET $2;`
	*/

	treesList = `
	SELECT id, name, description, is_active, create_time, creator, external_id
	FROM tree_mobile
	LIMIT $1 OFFSET $2;`
	/*
	treesGet = `
	SELECT 	thi.idthing AS id, 
			thi.name, 
			thi.description, 
			thi.isactive AS is_active, 
			thi.idcreator AS creator, 
			thi.datecreated AS create_time, 
			thi.dateinactivation AS inactivation_time, 
			thi.datelastmodif AS last_modification_time, 
			thi.idmodificator AS last_modification_user, 
			thi.isvalidated AS is_validated,
			null AS comment,
			null AS external_id,
			null AS id_validator,
			null AS inactivation_reason,
			'POINT(' || (pos.mineo/100.0)::varchar(10) || ' ' || (pos.minsn/100.0)::varchar(10) || ')' AS geom,
			jsonb_build_object('idvalidation', arbre.idvalidation, 'entouragerem', COALESCE(arbre.entouragerem,''), 'envracinairerem', COALESCE(arbre.envracinairerem,''), 'etatsanitairerem', COALESCE(arbre.etatsanitairerem,'')) as tree_attributes
	FROM thing thi
	INNER JOIN thing_position pos on pos.idthing = thi.idthing
	INNER JOIN thi_arbre arbre ON arbre.idthing = thi.idthing
	WHERE idtypething=74 AND thi.idthing=$1;`
	*/

	treesGet = `
	SELECT id, name, description, external_id, is_active, inactivation_time, inactivation_reason, comment, is_validated, id_validator,
			create_time, creator, last_modification_time, last_modification_user, geom, tree_attributes
	FROM tree_mobile
	WHERE id = $1`
	
	treesGetMaxId = "SELECT MAX(idthing) FROM thing WHERE idtypething=74"

	treesExist = "SELECT COUNT(*) FROM goeland_thing WHERE idthing=$1 AND idtypething=74;" 

	treesCount = "SELECT COUNT(*) FROM goeland_thing WHERE idtypething=74;"

	treesCreate = `
	INSERT INTO tree_mobile
	(name, description, external_id, is_active, comment, create_time, creator, geom, tree_attributes) 
	VALUES ($1, $2, $3, $4, $5, CURRENT_TIMESTAMP, $6, $7, $8)
	RETURNING id;`
	
	treesUpdate = ""

	treesDelete = ""

	treesSearchByName = "SELECT * FROM goeland_thing WHERE name LIKE $1 AND idtypething=74;"

	treesIsActive = "SELECT isactive FROM goeland_thing WHERE idthing=$1 AND idtypething=74;"

	treesCreateTable = `
	CREATE TABLE IF NOT EXISTS tree_mobile
	(
		id        				serial    			CONSTRAINT tree_mobile_pk   primary key,
		name					text	not null	constraint name_min_length check (length(btrim(name)) > 2),
		description				text				constraint description_min_length check (length(btrim(description)) > 2),
		external_id				int,
		is_active				boolean default true not null,
		inactivation_time		timestamp,
		inactivation_reason    	text,
		comment                	text,
		is_validated			boolean default false,
		id_validator			int,
		create_time				timestamp default now() not null,
		creator					integer	not null,
		last_modification_time	timestamp,
		last_modification_user	integer,
		geom					text	not null,
		tree_attributes			jsonb	not null
	);
	comment on table tree_mobile is 'tree_mobile is the main table of the sanarbo application';`
)