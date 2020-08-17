package postgres

import (
	"fmt"

	"github.com/google/uuid"
)

func getCorrectUUID(uuidStr string) (string, error) {
	newuuid, err := uuid.Parse(uuidStr)
	if err != nil {
		newuuid, err = uuid.NewUUID()
		if err != nil {
			return "", err
		}
	}

	return newuuid.String(), nil
}

func getCurrentLimmit(limit uint16) string {
	if limit > 0 {
		return fmt.Sprintf("LIMIT %d", limit)
	}
	return ""
}

func categoryAssocSubQuery(associate ...struct{}) (subQueryAssoc string) {
	if len(associate) > 0 {
		subQueryAssoc = `,
		COALESCE(
			(SELECT to_json(u) FROM users u
			WHERE pc.change_by_user = u.id
			AND u.deleted_at IS NULL), '{}') user_json,
		COALESCE(
			(SELECT json_agg(pr) FROM products pr
			WHERE pr.id IN (
				SELECT pcr.product_id FROM product_category_relations pcr
				WHERE pcr.category_id = pc.id
			) AND pr.deleted_at IS NULL), '{}') pr_json`
	}

	return subQueryAssoc
}

func productAssocSubQuery(associate ...struct{}) (subQueryAssoc string) {
	if len(associate) > 0 {
		subQueryAssoc = `,
		COALESCE(
			(SELECT to_json(u) FROM users u
			WHERE pr.change_by_user = u.id
				AND u.deleted_at IS NULL), '{}') user_json,
		COALESCE(
			(SELECT json_agg(pc) FROM product_categories pc
			WHERE pc.id IN (
					SELECT pcr.category_id FROM product_category_relations pcr
					WHERE pcr.product_id = pr.id
				) AND pc.deleted_at IS NULL), '{}') pr_json`
	}

	return subQueryAssoc
}
