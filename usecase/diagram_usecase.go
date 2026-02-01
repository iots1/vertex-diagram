package usecase

import (
	"context"
	"time"

	"github.com/iots1/vertex-diagram/domain"
)

type diagramUsecase struct {
	diagramRepo      domain.DiagramRepository
	tableRepo        domain.TableRepository
	relationshipRepo domain.RelationshipRepository
	contextTimeout   time.Duration
}

func NewDiagramUsecase(d domain.DiagramRepository, t domain.TableRepository, r domain.RelationshipRepository, timeout time.Duration) domain.DiagramUsecase {
	return &diagramUsecase{
		diagramRepo:      d,
		tableRepo:        t,
		relationshipRepo: r,
		contextTimeout:   timeout,
	}
}

func (u *diagramUsecase) GetAll(c context.Context) ([]domain.Diagram, error) {
	ctx, cancel := context.WithTimeout(c, u.contextTimeout)
	defer cancel()
	return u.diagramRepo.Fetch(ctx)
}

func (u *diagramUsecase) GetOne(c context.Context, id string) (*domain.Diagram, error) {
	ctx, cancel := context.WithTimeout(c, u.contextTimeout)
	defer cancel()

	// 1. Get diagram
	diagram, err := u.diagramRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	// 2. Get tables
	tables, err := u.tableRepo.GetByDiagramID(ctx, id)
	if err != nil {
		return nil, err
	}

	// 3. Get relationships
	relationships, err := u.relationshipRepo.GetByDiagramID(ctx, id)
	if err != nil {
		return nil, err
	}

	// 4. Merge tables and relationships back into content
	contentMap, ok := diagram.Content.(map[string]interface{})
	if !ok {
		contentMap = make(map[string]interface{})
	}

	contentMap["tables"] = tables
	contentMap["relationships"] = relationships
	diagram.Content = contentMap

	return diagram, nil
}

func (u *diagramUsecase) Save(c context.Context, d *domain.Diagram) (*domain.Diagram, error) {
	ctx, cancel := context.WithTimeout(c, u.contextTimeout)
	defer cancel()

	// 1. Save diagram first to get ID
	if d.ID == "" {
		if d.Name == "" {
			d.Name = "Untitled Diagram"
		}
		err := u.diagramRepo.Store(ctx, d)
		if err != nil {
			return nil, err
		}
	} else {
		err := u.diagramRepo.Update(ctx, d)
		if err != nil {
			return nil, err
		}
	}

	// 2. Clean up tables and relationships from content
	u.cleanupContent(d)

	// 3. Extract and save tables (now we have diagram ID)
	if err := u.saveTables(ctx, d); err != nil {
		return nil, err
	}

	// 4. Extract and save relationships (now we have diagram ID)
	if err := u.saveRelationships(ctx, d); err != nil {
		return nil, err
	}

	return d, nil
}

func (u *diagramUsecase) cleanupContent(d *domain.Diagram) {
	contentMap, ok := d.Content.(map[string]interface{})
	if !ok {
		return
	}

	delete(contentMap, "tables")
	delete(contentMap, "relationships")
	d.Content = contentMap
}

func (u *diagramUsecase) saveTables(ctx context.Context, d *domain.Diagram) error {
	contentMap, ok := d.Content.(map[string]interface{})
	if !ok {
		return nil
	}

	tablesData, ok := contentMap["tables"].([]interface{})
	if !ok {
		return nil
	}

	// Delete old tables first
	if err := u.tableRepo.DeleteByDiagramID(ctx, d.ID); err != nil {
		return err
	}

	// Extract and save new tables
	tables := make([]domain.Table, 0)
	for _, td := range tablesData {
		tableMap, ok := td.(map[string]interface{})
		if !ok {
			continue
		}

		table := domain.Table{
			DiagramID: d.ID,
			TableID:   getStringValue(tableMap, "id"),
			Name:      getStringValue(tableMap, "name"),
			Schema:    getStringValue(tableMap, "schema"),
			Fields:    tableMap["fields"],
			Indexes:   tableMap["indexes"],
			Color:     getStringValue(tableMap, "color"),
			X:         getIntValue(tableMap, "x"),
			Y:         getIntValue(tableMap, "y"),
			IsView:    getBoolValue(tableMap, "isView"),
			Order:     getIntValue(tableMap, "order"),
		}
		tables = append(tables, table)
	}

	if len(tables) > 0 {
		return u.tableRepo.StoreMultiple(ctx, tables)
	}
	return nil
}

func (u *diagramUsecase) saveRelationships(ctx context.Context, d *domain.Diagram) error {
	contentMap, ok := d.Content.(map[string]interface{})
	if !ok {
		return nil
	}

	relationshipsData, ok := contentMap["relationships"].([]interface{})
	if !ok {
		return nil
	}

	// Delete old relationships first
	if err := u.relationshipRepo.DeleteByDiagramID(ctx, d.ID); err != nil {
		return err
	}

	// Extract and save new relationships
	relationships := make([]domain.Relationship, 0)
	for _, rd := range relationshipsData {
		relMap, ok := rd.(map[string]interface{})
		if !ok {
			continue
		}

		rel := domain.Relationship{
			DiagramID:      d.ID,
			RelationshipID: getStringValue(relMap, "id"),
			Name:           getStringValue(relMap, "name"),
			SourceTableID:  getStringValue(relMap, "sourceTableId"),
			TargetTableID:  getStringValue(relMap, "targetTableId"),
			SourceFieldID:  getStringValue(relMap, "sourceFieldId"),
			TargetFieldID:  getStringValue(relMap, "targetFieldId"),
			Type:           getStringValue(relMap, "type"),
		}
		relationships = append(relationships, rel)
	}

	if len(relationships) > 0 {
		return u.relationshipRepo.StoreMultiple(ctx, relationships)
	}
	return nil
}

// Helper functions to safely extract values from maps
func getStringValue(m map[string]interface{}, key string) string {
	if v, ok := m[key]; ok {
		if s, ok := v.(string); ok {
			return s
		}
	}
	return ""
}

func getIntValue(m map[string]interface{}, key string) int {
	if v, ok := m[key]; ok {
		switch val := v.(type) {
		case float64:
			return int(val)
		case int:
			return val
		}
	}
	return 0
}

func getBoolValue(m map[string]interface{}, key string) bool {
	if v, ok := m[key]; ok {
		if b, ok := v.(bool); ok {
			return b
		}
	}
	return false
}

func (u *diagramUsecase) Delete(c context.Context, id string) error {
	ctx, cancel := context.WithTimeout(c, u.contextTimeout)
	defer cancel()

	// 1. Delete associated tables
	if err := u.tableRepo.DeleteByDiagramID(ctx, id); err != nil {
		return err
	}

	// 2. Delete associated relationships
	if err := u.relationshipRepo.DeleteByDiagramID(ctx, id); err != nil {
		return err
	}

	// 3. Delete the diagram
	return u.diagramRepo.Delete(ctx, id)
}