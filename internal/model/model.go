//nolint
//lint:file-ignore U1000 ignore unused code, it's generated
package model

import (
	"time"
)

type ColumnsEquipment struct {
	ID, Name, Description, Price string
}

type ColumnsEquipmentAssignment struct {
	ID, PersonID, EquipmentID, Deadline, CreatedAt string
	Equipment, Person                              string
}

type ColumnsEquipmentIncident struct {
	ID, EquipmentID, IncidentID, Deadline string
	Equipment, Incident                   string
}

type ColumnsIncident struct {
	ID, Ordinal, Description, CreatedAt, ResolvedAt, Deadline, AssigneeID, CreatorID, Status, Comment, Type, Priority, Approved string
	Assignee, Creator                                                                                                           string
}

type ColumnsPerson struct {
	ID, Login, Password, FullName, Email, Role, ManagerID string
	Manager                                               string
}

type ColumnsSchemaMigration struct {
	ID, Dirty string
}

type ColumnsSupport struct {
	ID, PersonID, IsManager string
	Person                  string
}

type ColumnsSt struct {
	Equipment           ColumnsEquipment
	EquipmentAssignment ColumnsEquipmentAssignment
	EquipmentIncident   ColumnsEquipmentIncident
	Incident            ColumnsIncident
	Person              ColumnsPerson
	SchemaMigration     ColumnsSchemaMigration
	Support             ColumnsSupport
}

var Columns = ColumnsSt{
	Equipment: ColumnsEquipment{
		ID:          "id",
		Name:        "name",
		Description: "description",
		Price:       "price",
	},
	EquipmentAssignment: ColumnsEquipmentAssignment{
		ID:          "id",
		PersonID:    "person_id",
		EquipmentID: "equipment_id",
		Deadline:    "deadline",
		CreatedAt:   "created_at",

		Equipment: "Equipment",
		Person:    "Person",
	},
	EquipmentIncident: ColumnsEquipmentIncident{
		ID:          "id",
		EquipmentID: "equipment_id",
		IncidentID:  "incident_id",
		Deadline:    "deadline",

		Equipment: "Equipment",
		Incident:  "Incident",
	},
	Incident: ColumnsIncident{
		ID:          "id",
		Ordinal:     "ordinal",
		Description: "description",
		CreatedAt:   "created_at",
		ResolvedAt:  "resolved_at",
		Deadline:    "deadline",
		AssigneeID:  "assignee_id",
		CreatorID:   "creator_id",
		Status:      "status",
		Comment:     "comment",
		Type:        "type",
		Priority:    "priority",
		Approved:    "approved",

		Assignee: "Assignee",
		Creator:  "Creator",
	},
	Person: ColumnsPerson{
		ID:        "id",
		Login:     "login",
		Password:  "password",
		FullName:  "full_name",
		Email:     "email",
		Role:      "role",
		ManagerID: "manager_id",

		Manager: "Manager",
	},
	SchemaMigration: ColumnsSchemaMigration{
		ID:    "version",
		Dirty: "dirty",
	},
	Support: ColumnsSupport{
		ID:        "id",
		PersonID:  "person_id",
		IsManager: "is_manager",

		Person: "Person",
	},
}

type TableEquipment struct {
	Name, Alias string
}

type TableEquipmentAssignment struct {
	Name, Alias string
}

type TableEquipmentIncident struct {
	Name, Alias string
}

type TableIncident struct {
	Name, Alias string
}

type TablePerson struct {
	Name, Alias string
}

type TableSchemaMigration struct {
	Name, Alias string
}

type TableSupport struct {
	Name, Alias string
}

type TablesSt struct {
	Equipment           TableEquipment
	EquipmentAssignment TableEquipmentAssignment
	EquipmentIncident   TableEquipmentIncident
	Incident            TableIncident
	Person              TablePerson
	SchemaMigration     TableSchemaMigration
	Support             TableSupport
}

var Tables = TablesSt{
	Equipment: TableEquipment{
		Name:  "equipment",
		Alias: "t",
	},
	EquipmentAssignment: TableEquipmentAssignment{
		Name:  "equipment_assignment",
		Alias: "t",
	},
	EquipmentIncident: TableEquipmentIncident{
		Name:  "equipment_incident",
		Alias: "t",
	},
	Incident: TableIncident{
		Name:  "incident",
		Alias: "t",
	},
	Person: TablePerson{
		Name:  "person",
		Alias: "t",
	},
	SchemaMigration: TableSchemaMigration{
		Name:  "schema_migrations",
		Alias: "t",
	},
	Support: TableSupport{
		Name:  "support",
		Alias: "t",
	},
}

type Equipment struct {
	tableName struct{} `sql:"equipment,alias:t" pg:",discard_unknown_columns"`

	ID          string `sql:"id,pk,type:uuid"`
	Name        string `sql:"name,notnull"`
	Description string `sql:"description,notnull"`
	Price       int    `sql:"price,notnull"`
}

type EquipmentAssignment struct {
	tableName struct{} `sql:"equipment_assignment,alias:t" pg:",discard_unknown_columns"`

	ID          string    `sql:"id,pk,type:uuid"`
	PersonID    string    `sql:"person_id,type:uuid,notnull"`
	EquipmentID string    `sql:"equipment_id,type:uuid,notnull"`
	Deadline    time.Time `sql:"deadline,notnull"`
	CreatedAt   time.Time `sql:"created_at,notnull"`

	Equipment *Equipment `pg:"fk:equipment_id"`
	Person    *Person    `pg:"fk:person_id"`
}

type EquipmentIncident struct {
	tableName struct{} `sql:"equipment_incident,alias:t" pg:",discard_unknown_columns"`

	ID          string    `sql:"id,pk,type:uuid"`
	EquipmentID string    `sql:"equipment_id,type:uuid,notnull"`
	IncidentID  string    `sql:"incident_id,type:uuid,notnull"`
	Deadline    time.Time `sql:"deadline,notnull"`

	Equipment *Equipment `pg:"fk:equipment_id"`
	Incident  *Incident  `pg:"fk:incident_id"`
}

type Incident struct {
	tableName struct{} `sql:"incident,alias:t" pg:",discard_unknown_columns"`

	ID          string     `sql:"id,pk,type:uuid"`
	Ordinal     *int       `sql:"ordinal"`
	Description string     `sql:"description,notnull"`
	CreatedAt   time.Time  `sql:"created_at,notnull"`
	ResolvedAt  *time.Time `sql:"resolved_at"`
	Deadline    time.Time  `sql:"deadline,notnull"`
	AssigneeID  *string    `sql:"assignee_id,type:uuid"`
	CreatorID   string     `sql:"creator_id,type:uuid,notnull"`
	Status      string     `sql:"status,notnull"`
	Comment     *string    `sql:"comment"`
	Type        string     `sql:"type,notnull"`
	Priority    int        `sql:"priority,notnull"`
	Approved    bool       `sql:"approved,notnull"`

	Assignee *Support `pg:"fk:assignee_id"`
	Creator  *Person  `pg:"fk:creator_id"`
}

type Person struct {
	tableName struct{} `sql:"person,alias:t" pg:",discard_unknown_columns"`

	ID        string  `sql:"id,pk,type:uuid"`
	Login     string  `sql:"login,notnull"`
	Password  string  `sql:"password,notnull"`
	FullName  string  `sql:"full_name,notnull"`
	Email     string  `sql:"email,notnull"`
	Role      string  `sql:"role,notnull"`
	ManagerID *string `sql:"manager_id,type:uuid"`

	Manager *Person `pg:"fk:manager_id"`
}

type SchemaMigration struct {
	tableName struct{} `sql:"schema_migrations,alias:t" pg:",discard_unknown_columns"`

	ID    int64 `sql:"version,pk"`
	Dirty bool  `sql:"dirty,notnull"`
}

type Support struct {
	tableName struct{} `sql:"support,alias:t" pg:",discard_unknown_columns"`

	ID        string `sql:"id,pk,type:uuid"`
	PersonID  string `sql:"person_id,type:uuid,notnull"`
	IsManager bool   `sql:"is_manager,notnull"`

	Person *Person `pg:"fk:person_id"`
}
