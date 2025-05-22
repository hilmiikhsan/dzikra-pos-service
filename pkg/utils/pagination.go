package utils

func Paginate(page, limit int) (currentPage, perPage, offset int) {
	if limit <= 0 {
		limit = 10
	}

	if page <= 0 {
		page = 1
	}

	currentPage = page
	perPage = limit
	offset = (page - 1) * limit
	return
}

func CalculateTotalPages(total, perPage int) int {
	if perPage <= 0 {
		perPage = 10
	}

	if total == 0 {
		return 1
	}

	return (total + perPage - 1) / perPage
}
