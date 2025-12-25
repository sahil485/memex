package client

func ConfigureIndexSettings() error {
	c := New()
	index := c.GetIndex()

	_, err := index.UpdateSearchableAttributes(&[]string{
		"name",
		"title",
		"content",
		"tags",
		"path",
	})
	if err != nil {
		return err
	}

	_, err = index.UpdateFilterableAttributes(&[]interface{}{
		"ext",
		"dir",
		"mod_time",
		"size",
		"tags",
	})
	if err != nil {
		return err
	}

	_, err = index.UpdateSortableAttributes(&[]string{
		"mod_time",
		"size",
		"name",
	})
	if err != nil {
		return err
	}

	_, err = index.UpdateDisplayedAttributes(&[]string{
		"id",
		"path",
		"name",
		"dir",
		"ext",
		"size",
		"mod_time",
		"title",
		"tags",
	})

	return err
}
