package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/field"
	"github.com/google/uuid"
	"huspass/model"
)

// Msisdn holds the schema definition for the Msisdn entity.
type Msisdn struct {
	ent.Schema
}

// Fields of the Msisdn.
func (Msisdn) Fields() []ent.Field {
	return []ent.Field{
		field.UUID("id", uuid.UUID{}).Default(uuid.New),
		field.String("primary_identity").NotEmpty(),
		field.Enum("primary_identity_type").GoType(model.PrimaryIdentityType("")),
		field.Bool("provisioned").Default(false),
	}
}

// Edges of the Msisdn.
func (Msisdn) Edges() []ent.Edge {
	return nil
}
