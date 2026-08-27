package amazontextract

import (
	"strings"

	"aidoc/internal/entity"

	"github.com/aws/aws-sdk-go-v2/service/textract/types"
)

func ExtractStructuredData(doc *entity.OCRDocument, blocks []types.Block, features []entity.TextractFeature) {
	for _, f := range features {
		switch f {
		case entity.TextractForms:
			ExtractKeyValuePairs(doc, blocks)
		case entity.TextractTables:
			ExtractTables(doc, blocks)
		case entity.TextractLayout:
			ExtractLayout(doc, blocks)
		case entity.TextractQueries:
			ExtractQueryResults(doc, blocks)
		}
	}
}

func ExtractKeyValuePairs(doc *entity.OCRDocument, blocks []types.Block) {
	blockMap := BuildBlockMap(blocks)

	fields := make(map[string]string)
	for _, b := range blocks {
		if b.BlockType != types.BlockTypeKeyValueSet || b.EntityTypes == nil {
			continue
		}

		isKey := false
		for _, et := range b.EntityTypes {
			if et == types.EntityTypeKey {
				isKey = true
				break
			}
		}
		if !isKey || b.Relationships == nil {
			continue
		}

		keyText := ExtractTextFromRelationships(blockMap, b.Relationships, types.RelationshipTypeChild)
		if keyText == "" {
			continue
		}

		var valueText string
		for _, rel := range b.Relationships {
			if rel.Type != types.RelationshipTypeValue {
				continue
			}
			for _, id := range rel.Ids {
				if valueBlock, ok := blockMap[id]; ok {
					valueText = ExtractTextFromRelationships(blockMap, valueBlock.Relationships, types.RelationshipTypeChild)
					if valueText == "" {
						valueText = ExtractSelectionElement(blockMap, valueBlock.Relationships)
					}
				}
			}
		}

		if valueText != "" {
			fields[keyText] = valueText
		}
	}

	if len(fields) > 0 {
		data := make(map[string]interface{})
		for k, v := range fields {
			data[k] = v
		}
		doc.ExtractedData["form_fields"] = data
	}
}

func ExtractSelectionElement(blockMap map[string]types.Block, relationships []types.Relationship) string {
	for _, rel := range relationships {
		if rel.Type != types.RelationshipTypeChild {
			continue
		}
		for _, id := range rel.Ids {
			if child, ok := blockMap[id]; ok && child.BlockType == types.BlockTypeSelectionElement {
				if child.SelectionStatus == types.SelectionStatusSelected {
					return "selected"
				}
				return "not_selected"
			}
		}
	}
	return ""
}

func ExtractTables(doc *entity.OCRDocument, blocks []types.Block) {
	blockMap := BuildBlockMap(blocks)

	var tables []map[string]interface{}
	for _, b := range blocks {
		if b.BlockType != types.BlockTypeTable {
			continue
		}

		var tableTitle, tableFooter string
		for _, rel := range b.Relationships {
			if rel.Type != types.RelationshipTypeChild {
				continue
			}
			for _, id := range rel.Ids {
				if child, ok := blockMap[id]; ok {
					switch child.BlockType {
					case types.BlockTypeTableTitle:
						tableTitle = ExtractTextFromRelationships(blockMap, child.Relationships, types.RelationshipTypeChild)
					case types.BlockTypeTableFooter:
						tableFooter = ExtractTextFromRelationships(blockMap, child.Relationships, types.RelationshipTypeChild)
					}
				}
			}
		}

		rows := make(map[int]map[int]string)
		for _, rel := range b.Relationships {
			if rel.Type != types.RelationshipTypeChild {
				continue
			}
			for _, id := range rel.Ids {
				if cell, ok := blockMap[id]; ok && cell.BlockType == types.BlockTypeCell {
					row := int(*cell.RowIndex) - 1
					col := int(*cell.ColumnIndex) - 1
					text := ExtractTextFromRelationships(blockMap, cell.Relationships, types.RelationshipTypeChild)
					if rows[row] == nil {
						rows[row] = make(map[int]string)
					}
					rows[row][col] = text
				}
			}
		}

		if len(rows) > 0 {
			var tableData []map[string]string
			headerRow := rows[0]
			for rowIdx := 1; rowIdx <= len(rows); rowIdx++ {
				rowData := rows[rowIdx-1]
				if rowData == nil {
					continue
				}
				row := make(map[string]string)
				for colIdx := 0; colIdx < len(rowData); colIdx++ {
					header := headerRow[colIdx]
					if header == "" {
						header = "column_0"
					}
					row[header] = rowData[colIdx]
				}
				tableData = append(tableData, row)
			}
			tables = append(tables, map[string]interface{}{
				"title":  tableTitle,
				"footer": tableFooter,
				"header": headerRow,
				"rows":   tableData,
			})
		}
	}

	if len(tables) > 0 {
		doc.ExtractedData["tables"] = tables
	}
}

func ExtractLayout(doc *entity.OCRDocument, blocks []types.Block) {
	layoutElements := make(map[string][]string)
	for _, block := range blocks {
		switch block.BlockType {
		case types.BlockTypeLayoutHeader,
			types.BlockTypeLayoutFooter,
			types.BlockTypeLayoutSectionHeader,
			types.BlockTypeLayoutPageNumber,
			types.BlockTypeLayoutList,
			types.BlockTypeLayoutFigure,
			types.BlockTypeLayoutKeyValue,
			types.BlockTypeLayoutTitle,
			types.BlockTypeLayoutText,
			types.BlockTypeLayoutTable:
			key := string(block.BlockType)
			if block.Text != nil {
				layoutElements[key] = append(layoutElements[key], *block.Text)
			}
		}
	}
	if len(layoutElements) > 0 {
		doc.ExtractedData["layout"] = layoutElements
	}
}

func ExtractQueryResults(doc *entity.OCRDocument, blocks []types.Block) {
	queryResults := make([]map[string]interface{}, 0)
	for _, block := range blocks {
		if block.BlockType == types.BlockTypeQueryResult {
			result := make(map[string]interface{})
			if block.Query != nil {
				if block.Query.Alias != nil {
					result["alias"] = *block.Query.Alias
				}
				if block.Query.Text != nil {
					result["query"] = *block.Query.Text
				}
			}
			if block.Text != nil {
				result["answer"] = *block.Text
			}
			if block.Confidence != nil {
				result["confidence"] = float64(*block.Confidence)
			}
			if block.Page != nil {
				result["page"] = int(*block.Page)
			}
			queryResults = append(queryResults, result)
		}
	}
	if len(queryResults) > 0 {
		doc.ExtractedData["query_results"] = queryResults
	}
}

func ExtractTextFromRelationships(blockMap map[string]types.Block, relationships []types.Relationship, relType types.RelationshipType) string {
	var parts []string
	for _, rel := range relationships {
		if rel.Type != relType {
			continue
		}
		for _, id := range rel.Ids {
			if child, ok := blockMap[id]; ok {
				if child.Text != nil {
					parts = append(parts, *child.Text)
				}
			}
		}
	}
	return strings.Join(parts, " ")
}

func BuildBlockMap(blocks []types.Block) map[string]types.Block {
	blockMap := make(map[string]types.Block)
	for _, b := range blocks {
		blockMap[*b.Id] = b
	}
	return blockMap
}

func BuildFeatureTypes(features []entity.TextractFeature) []types.FeatureType {
	featureTypes := make([]types.FeatureType, 0, len(features))
	for _, f := range features {
		switch f {
		case entity.TextractForms:
			featureTypes = append(featureTypes, types.FeatureTypeForms)
		case entity.TextractTables:
			featureTypes = append(featureTypes, types.FeatureTypeTables)
		case entity.TextractSignature:
			featureTypes = append(featureTypes, types.FeatureTypeSignatures)
		case entity.TextractLayout, entity.TextractID:
			featureTypes = append(featureTypes, types.FeatureTypeLayout)
		case entity.TextractQueries:
			featureTypes = append(featureTypes, types.FeatureTypeQueries)
		}
	}
	return featureTypes
}

func BuildQueriesConfig(queries []entity.Query) *types.QueriesConfig {
	qs := make([]types.Query, len(queries))
	for i, q := range queries {
		qs[i] = types.Query{
			Text:  &q.Text,
			Alias: &q.Alias,
		}
	}
	return &types.QueriesConfig{
		Queries: qs,
	}
}
