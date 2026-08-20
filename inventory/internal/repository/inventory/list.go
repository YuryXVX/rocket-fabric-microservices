package inventory

import (
	"context"
	"slices"
	"strings"

	"github.com/google/uuid"
	errs "inventory/internal/errors"
	"inventory/internal/model"
	"inventory/internal/repository/converter"
	"inventory/internal/repository/record"
	"inventory/internal/service/input"
)

func (r *repository) List(ctx context.Context, input input.PartFilter) ([]*model.Part, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	parts := filterByUUID(r.parts, input)

	if len(input.UUIDs) > 0 && len(input.UUIDs) != len(parts) {
		return []*model.Part{}, errs.ErrPartNotFound
	}

	return parts, nil
}

func filterByUUID(parts map[string]record.Part, input input.PartFilter) []*model.Part {
	filtered := make([]*model.Part, 0, len(parts))

	if len(input.UUIDs) > 0 {
		for _, p := range parts {
			if slices.Contains(input.UUIDs, uuid.MustParse(p.UUID)) {
				filtered = append(filtered, converter.RecordPartToModel(p))
			}
		}

		return filtered
	} else if input.PartType != model.PartTypeUnspecified {
		for _, p := range parts {
			if p.PartType == converter.ConvertRecordPartType(input.PartType) {
				filtered = append(filtered, converter.RecordPartToModel(p))
			}
		}

		sortByName(&filtered)

		return filtered
	}

	for _, p := range parts {
		filtered = append(filtered, converter.RecordPartToModel(p))
	}

	sortByName(&filtered)

	return filtered
}

func sortByName(parts *[]*model.Part) {
	slices.SortFunc(*parts, func(a, b *model.Part) int {
		if a == nil && b == nil {
			return 0
		}
		if a == nil {
			return -1
		}
		if b == nil {
			return 1
		}

		nameA := strings.ToLower(a.Name)
		nameB := strings.ToLower(b.Name)

		if nameA < nameB {
			return -1
		}
		if nameA > nameB {
			return 1
		}
		return 0
	})
}
