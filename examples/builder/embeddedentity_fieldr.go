// Code generated by 'fieldr'; DO NOT EDIT.

package builder

type EmbeddedEntityBuilder struct {
	Metadata string
}

func NewEmbeddedEntityBuilder() *EmbeddedEntityBuilder {
	return &EmbeddedEntityBuilder{}
}

func (b EmbeddedEntityBuilder) Build() EmbeddedEntity {
	return EmbeddedEntity{
		Metadata: b.Metadata,
	}
}

func (b EmbeddedEntityBuilder) SetMetadata(metadata string) EmbeddedEntityBuilder {
	b.Metadata = metadata
	return b
}

func (i EmbeddedEntity) ToBuilder() EmbeddedEntityBuilder {
	return EmbeddedEntityBuilder{
		Metadata: i.Metadata,
	}
}
