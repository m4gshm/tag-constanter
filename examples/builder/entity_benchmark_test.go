package builder

import (
	"testing"
)

func BenchmarkNoBuilder(b *testing.B) {
	var e Entity[int32]
	for i := 0; i < b.N; i++ {
		e = Entity[int32]{
			BaseEntity: &BaseEntity[int32]{
				ID:                   2,
				CodeAwareEntity:      &CodeAwareEntity{},
				ForeignIDAwareEntity: ForeignIDAwareEntity[int32]{},
			},
			Name:     "1",
			Surname:  "3",
			Embedded: EmbeddedEntity{Metadata: "meta"},
		}
	}
	_ = e
}

func BenchmarkBuilderValNoChain(b *testing.B) {
	var e Entity[int32]
	for i := 0; i < b.N; i++ {
		e = EntityBuilderVal[int32]{Name: "1", ID: 2, Surname: "3", Embedded: EmbeddedEntityBuilder{Metadata: "meta"}.Build()}.Build()
	}
	_ = e
}

func BenchmarkBuilderRefToVal(b *testing.B) {
	var e Entity[int32]
	for i := 0; i < b.N; i++ {
		e = *NewEntityBuilder[int32]().SetName("1").SetID(2).SetSurname("3").SetEmbedded(EmbeddedEntityBuilder{}.SetMetadata("meta").Build()).Build()
	}
	_ = e
}

func BenchmarkBuilderChainRefBuildVal(b *testing.B) {
	var e Entity[int32]
	for i := 0; i < b.N; i++ {
		e = (&EntityBuilderChainRefBuildVal[int32]{Name: "1"}).SetID(2).SetSurname("3").SetEmbedded(EmbeddedEntityBuilder{}.SetMetadata("meta").Build()).Build()
	}
	_ = e
}

func BenchmarkBuilderVal(b *testing.B) {
	var e Entity[int32]
	for i := 0; i < b.N; i++ {
		e = EntityBuilderVal[int32]{Name: "1"}.SetID(2).SetSurname("3").SetEmbedded(EmbeddedEntityBuilder{}.SetMetadata("meta").Build()).Build()
	}
	_ = e
}

func BenchmarkBuilderRef(b *testing.B) {
	var e *Entity[int32]
	for i := 0; i < b.N; i++ {
		e = (&EntityBuilder[int32]{}).SetName("1").SetID(2).SetSurname("3").SetEmbedded(EmbeddedEntityBuilder{}.SetMetadata("meta").Build()).Build()
	}
	_ = e
}

func BenchmarkBuilderValToRef(b *testing.B) {
	var e *Entity[int32]
	for i := 0; i < b.N; i++ {
		v := EntityBuilderVal[int32]{Name: "1"}.SetID(2).SetSurname("3").SetEmbedded(EmbeddedEntityBuilder{}.SetMetadata("meta").Build()).Build()
		e = &v
	}
	_ = e
}
