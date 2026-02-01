package usecase

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/iots1/vertex-diagram/domain"
)

type diagramUsecase struct {
	diagramRepo       domain.DiagramRepository
	tableRepo         domain.TableRepository
	relationshipRepo  domain.RelationshipRepository
	dependencyRepo    domain.DependencyRepository
	areaRepo          domain.AreaRepository
	customTypeRepo    domain.CustomTypeRepository
	noteRepo          domain.NoteRepository
	diagramFilterRepo domain.DiagramFilterRepository
	contextTimeout    time.Duration
}

func NewDiagramUsecase(
	d domain.DiagramRepository,
	t domain.TableRepository,
	r domain.RelationshipRepository,
	dep domain.DependencyRepository,
	area domain.AreaRepository,
	ct domain.CustomTypeRepository,
	note domain.NoteRepository,
	df domain.DiagramFilterRepository,
	timeout time.Duration,
) domain.DiagramUsecase {
	return &diagramUsecase{
		diagramRepo:       d,
		tableRepo:         t,
		relationshipRepo:  r,
		dependencyRepo:    dep,
		areaRepo:          area,
		customTypeRepo:    ct,
		noteRepo:          note,
		diagramFilterRepo: df,
		contextTimeout:    timeout,
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

	// 4. Get dependencies
	dependencies, err := u.dependencyRepo.GetByDiagramID(ctx, id)
	if err != nil {
		return nil, err
	}

	// 5. Get areas
	areas, err := u.areaRepo.GetByDiagramID(ctx, id)
	if err != nil {
		return nil, err
	}

	// 6. Get custom types
	customTypes, err := u.customTypeRepo.GetByDiagramID(ctx, id)
	if err != nil {
		return nil, err
	}

	// 7. Get notes
	notes, err := u.noteRepo.GetByDiagramID(ctx, id)
	if err != nil {
		return nil, err
	}

	// 8. Get diagram filter
	diagramFilter, err := u.diagramFilterRepo.GetByDiagramID(ctx, id)
	if err != nil {
		return nil, err
	}

	// 9. Merge all entities back into content
	if diagram.Content == nil {
		diagram.Content = make(map[string]interface{})
	}

	diagram.Content["tables"] = tables
	diagram.Content["relationships"] = relationships
	diagram.Content["dependencies"] = dependencies
	diagram.Content["areas"] = areas
	diagram.Content["customTypes"] = customTypes
	diagram.Content["notes"] = notes
	if diagramFilter != nil {
		diagram.Content["diagramFilter"] = diagramFilter
	}

	return diagram, nil
}

func (u *diagramUsecase) Save(c context.Context, d *domain.Diagram) (*domain.Diagram, error) {
	ctx, cancel := context.WithTimeout(c, u.contextTimeout)
	defer cancel()

	if d == nil {
		return nil, fmt.Errorf("diagram is nil")
	}

	log.Printf("💾 Saving diagram: ID=%s, Name=%s", d.ID, d.Name)

	// 1. Save diagram first to get ID
	if d.ID == "" {
		if d.Name == "" {
			d.Name = "Untitled Diagram"
		}
		log.Printf("  📌 Creating new diagram: %s", d.Name)
		err := u.diagramRepo.Store(ctx, d)
		if err != nil {
			log.Printf("  ❌ Error storing diagram: %v", err)
			return nil, err
		}
	} else {
		log.Printf("  📌 Updating existing diagram: ID=%s", d.ID)
		err := u.diagramRepo.Update(ctx, d)
		if err != nil {
			log.Printf("  ❌ Error updating diagram: %v", err)
			return nil, err
		}
	}

	// 2. Extract and save all entities BEFORE cleaning up content
	log.Printf("  📋 Saving tables...")
	if err := u.saveTables(ctx, d); err != nil {
		log.Printf("  ❌ Error saving tables: %v", err)
		return nil, err
	}

	log.Printf("  🔗 Saving relationships...")
	if err := u.saveRelationships(ctx, d); err != nil {
		log.Printf("  ❌ Error saving relationships: %v", err)
		return nil, err
	}

	log.Printf("  ⛓️  Saving dependencies...")
	if err := u.saveDependencies(ctx, d); err != nil {
		log.Printf("  ❌ Error saving dependencies: %v", err)
		return nil, err
	}

	log.Printf("  📦 Saving areas...")
	if err := u.saveAreas(ctx, d); err != nil {
		log.Printf("  ❌ Error saving areas: %v", err)
		return nil, err
	}

	log.Printf("  🎨 Saving custom types...")
	if err := u.saveCustomTypes(ctx, d); err != nil {
		log.Printf("  ❌ Error saving custom types: %v", err)
		return nil, err
	}

	log.Printf("  📝 Saving notes...")
	if err := u.saveNotes(ctx, d); err != nil {
		log.Printf("  ❌ Error saving notes: %v", err)
		return nil, err
	}

	log.Printf("  🔍 Saving diagram filter...")
	if err := u.saveDiagramFilter(ctx, d); err != nil {
		log.Printf("  ❌ Error saving diagram filter: %v", err)
		return nil, err
	}

	// 3. Clean up entity arrays from content AFTER extracting
	log.Printf("  🧹 Cleaning up content...")
	u.cleanupContent(d)

	// 4. Update diagram with cleaned content (no entity arrays)
	log.Printf("  💾 Saving cleaned diagram content...")
	err := u.diagramRepo.Update(ctx, d)
	if err != nil {
		log.Printf("  ❌ Error saving cleaned diagram: %v", err)
		return nil, err
	}

	log.Printf("✅ Diagram saved successfully: ID=%s", d.ID)
	return d, nil
}

func (u *diagramUsecase) cleanupContent(d *domain.Diagram) {
	if d.Content == nil {
		return
	}

	delete(d.Content, "tables")
	delete(d.Content, "relationships")
	delete(d.Content, "dependencies")
	delete(d.Content, "areas")
	delete(d.Content, "customTypes")
	delete(d.Content, "notes")
	delete(d.Content, "diagramFilter")
}

func (u *diagramUsecase) saveTables(ctx context.Context, d *domain.Diagram) error {
	if d.Content == nil {
		return nil
	}

	tablesData, ok := d.Content["tables"].([]interface{})
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

		// Convert fields to []map[string]interface{}
		fields := getMapArrayValue(tableMap, "fields")
		indexes := getMapArrayValue(tableMap, "indexes")

		table := domain.Table{
			DiagramID: d.ID,
			TableID:   getStringValue(tableMap, "id"),
			Name:      getStringValue(tableMap, "name"),
			Schema:    getStringValue(tableMap, "schema"),
			Fields:    fields,
			Indexes:   indexes,
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
	if d.Content == nil {
		return nil
	}

	relationshipsData, ok := d.Content["relationships"].([]interface{})
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

func (u *diagramUsecase) saveDependencies(ctx context.Context, d *domain.Diagram) error {
	if d.Content == nil {
		return nil
	}

	dependenciesData, ok := d.Content["dependencies"].([]interface{})
	if !ok {
		return nil
	}

	// Delete old dependencies first
	if err := u.dependencyRepo.DeleteByDiagramID(ctx, d.ID); err != nil {
		return err
	}

	// Extract and save new dependencies
	dependencies := make([]domain.Dependency, 0)
	for _, dd := range dependenciesData {
		depMap, ok := dd.(map[string]interface{})
		if !ok {
			continue
		}

		dep := domain.Dependency{
			DiagramID:        d.ID,
			DependencyID:     getStringValue(depMap, "id"),
			Schema:           getStringValue(depMap, "schema"),
			TableID:          getStringValue(depMap, "tableId"),
			DependentSchema:  getStringValue(depMap, "dependentSchema"),
			DependentTableID: getStringValue(depMap, "dependentTableId"),
		}
		dependencies = append(dependencies, dep)
	}

	if len(dependencies) > 0 {
		return u.dependencyRepo.StoreMultiple(ctx, dependencies)
	}
	return nil
}

func (u *diagramUsecase) saveAreas(ctx context.Context, d *domain.Diagram) error {
	if d.Content == nil {
		return nil
	}

	areasData, ok := d.Content["areas"].([]interface{})
	if !ok {
		return nil
	}

	// Delete old areas first
	if err := u.areaRepo.DeleteByDiagramID(ctx, d.ID); err != nil {
		return err
	}

	// Extract and save new areas
	areas := make([]domain.Area, 0)
	for _, ad := range areasData {
		areaMap, ok := ad.(map[string]interface{})
		if !ok {
			continue
		}

		area := domain.Area{
			DiagramID: d.ID,
			Name:      getStringValue(areaMap, "name"),
			X:         getIntValue(areaMap, "x"),
			Y:         getIntValue(areaMap, "y"),
			Width:     getIntValue(areaMap, "width"),
			Height:    getIntValue(areaMap, "height"),
			Color:     getStringValue(areaMap, "color"),
		}
		areas = append(areas, area)
	}

	if len(areas) > 0 {
		return u.areaRepo.StoreMultiple(ctx, areas)
	}
	return nil
}

func (u *diagramUsecase) saveCustomTypes(ctx context.Context, d *domain.Diagram) error {
	if d.Content == nil {
		return nil
	}

	customTypesData, ok := d.Content["customTypes"].([]interface{})
	if !ok {
		return nil
	}

	// Delete old custom types first
	if err := u.customTypeRepo.DeleteByDiagramID(ctx, d.ID); err != nil {
		return err
	}

	// Extract and save new custom types
	customTypes := make([]domain.CustomType, 0)
	for _, ctd := range customTypesData {
		ctMap, ok := ctd.(map[string]interface{})
		if !ok {
			continue
		}

		ct := domain.CustomType{
			DiagramID: d.ID,
			Schema:    getStringValue(ctMap, "schema"),
			Type:      getStringValue(ctMap, "type"),
			Kind:      getStringValue(ctMap, "kind"),
			Values:    ctMap["values"],
			Fields:    ctMap["fields"],
		}
		customTypes = append(customTypes, ct)
	}

	if len(customTypes) > 0 {
		return u.customTypeRepo.StoreMultiple(ctx, customTypes)
	}
	return nil
}

func (u *diagramUsecase) saveNotes(ctx context.Context, d *domain.Diagram) error {
	if d.Content == nil {
		return nil
	}

	notesData, ok := d.Content["notes"].([]interface{})
	if !ok {
		return nil
	}

	// Delete old notes first
	if err := u.noteRepo.DeleteByDiagramID(ctx, d.ID); err != nil {
		return err
	}

	// Extract and save new notes
	notes := make([]domain.Note, 0)
	for _, nd := range notesData {
		noteMap, ok := nd.(map[string]interface{})
		if !ok {
			continue
		}

		note := domain.Note{
			DiagramID: d.ID,
			Content:   getStringValue(noteMap, "content"),
			X:         getIntValue(noteMap, "x"),
			Y:         getIntValue(noteMap, "y"),
			Width:     getIntValue(noteMap, "width"),
			Height:    getIntValue(noteMap, "height"),
			Color:     getStringValue(noteMap, "color"),
		}
		notes = append(notes, note)
	}

	if len(notes) > 0 {
		return u.noteRepo.StoreMultiple(ctx, notes)
	}
	return nil
}

func (u *diagramUsecase) saveDiagramFilter(ctx context.Context, d *domain.Diagram) error {
	if d == nil || d.ID == "" {
		return nil
	}

	if d.Content == nil {
		return nil
	}

	filterData, ok := d.Content["diagramFilter"].(map[string]interface{})
	if !ok {
		// No filter data, optionally delete existing filter
		return u.diagramFilterRepo.DeleteByDiagramID(ctx, d.ID)
	}

	// Extract table IDs
	tableIDs := make([]string, 0)
	if tableIDsData, ok := filterData["tableIds"].([]interface{}); ok {
		for _, id := range tableIDsData {
			if str, ok := id.(string); ok {
				tableIDs = append(tableIDs, str)
			}
		}
	}

	// Extract schema IDs
	schemaIDs := make([]string, 0)
	if schemaIDsData, ok := filterData["schemaIds"].([]interface{}); ok {
		for _, id := range schemaIDsData {
			if str, ok := id.(string); ok {
				schemaIDs = append(schemaIDs, str)
			}
		}
	}

	// If no actual filter data, delete the filter
	if len(tableIDs) == 0 && len(schemaIDs) == 0 {
		return u.diagramFilterRepo.DeleteByDiagramID(ctx, d.ID)
	}

	filter := domain.DiagramFilter{
		DiagramID: d.ID,
		TableIDs:  tableIDs,
		SchemaIDs: schemaIDs,
	}

	return u.diagramFilterRepo.Store(ctx, &filter)
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

func getMapArrayValue(m map[string]interface{}, key string) []map[string]interface{} {
	if v, ok := m[key]; ok {
		// Handle []interface{}
		if arrInterface, ok := v.([]interface{}); ok {
			result := make([]map[string]interface{}, 0)
			for _, item := range arrInterface {
				if itemMap, ok := item.(map[string]interface{}); ok {
					result = append(result, itemMap)
				}
			}
			return result
		}
		// Handle []map[string]interface{} directly
		if arrMap, ok := v.([]map[string]interface{}); ok {
			return arrMap
		}
	}
	return []map[string]interface{}{}
}

func (u *diagramUsecase) Delete(c context.Context, id string) error {
	ctx, cancel := context.WithTimeout(c, u.contextTimeout)
	defer cancel()

	// Cascade delete all associated entities
	if err := u.tableRepo.DeleteByDiagramID(ctx, id); err != nil {
		return err
	}

	if err := u.relationshipRepo.DeleteByDiagramID(ctx, id); err != nil {
		return err
	}

	if err := u.dependencyRepo.DeleteByDiagramID(ctx, id); err != nil {
		return err
	}

	if err := u.areaRepo.DeleteByDiagramID(ctx, id); err != nil {
		return err
	}

	if err := u.customTypeRepo.DeleteByDiagramID(ctx, id); err != nil {
		return err
	}

	if err := u.noteRepo.DeleteByDiagramID(ctx, id); err != nil {
		return err
	}

	if err := u.diagramFilterRepo.DeleteByDiagramID(ctx, id); err != nil {
		return err
	}

	// Finally, delete the diagram
	return u.diagramRepo.Delete(ctx, id)
}