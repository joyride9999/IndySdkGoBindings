package inMemUtils

import "github.com/Jeffail/gabs/v2"

// Duplicate from unitTestUtils to overcome import cycle.

func IsIncluded(expected *gabs.Container, resulted *gabs.Container) bool {
	ok := false

	// Iterate through expected items.
	for path, element := range expected.ChildrenMap() {
		// Search elements by their path in resulted.
		search := resulted.Path(path)
		switch search.Data().(type) {
		case []interface{}:
			// Element is an array.
			ok = true
			if len(element.Children()) != len(search.Children()) {
				return false
			}
			// Iterate through array's components.
			for _, component := range element.Children() {
				exists := false
				for _, item := range search.Children() {
					exists = isEqual(component, item)
					if exists {
						break
					}
				}
				if !exists {
					ok = false
					break
				}
			}

		case map[string]interface{}:
			// Exception if *gabs.Container is found, but it's empty.
			if len(search.ChildrenMap()) == 0 {
				ok = true
			} else {
				// Element is a *gabs.Container.
				ok = IsIncluded(element, search)
			}
		case nil:
			// Element is not found in resulted JSON.
			ok = false
		default:
			if search.Data() == element.Data() {
				ok = true
			} else {
				ok = false
			}
		}
		if ok == false {
			break
		}
	}
	return ok
}

func isEqual(expected *gabs.Container, resulted *gabs.Container) bool {
	ok := false

	if len(expected.ChildrenMap()) != len(resulted.ChildrenMap()) {
		return ok
	} else {
		// If expected underlying value is an object, the map of children isn't empty.
		if len(expected.ChildrenMap()) != 0 && len(resulted.ChildrenMap()) != 0 {
			// Iterate through expected items.
			for path, element := range expected.ChildrenMap() {
				// Search elements by their path in resulted.
				search := resulted.Path(path)
				switch v := search.Data().(type) {
				case []interface{}:
					// Element is an array.
					ok = true
					if len(element.Children()) != len(search.Children()) {
						return false
					}
					// Iterate through array's components.
					for _, component := range element.Children() {
						exists := false
						for _, item := range search.Children() {
							exists = isEqual(component, item)
							if exists {
								break
							}
						}
						if !exists {
							ok = false
							break
						}
					}

				case map[string]interface{}:
					// Element is a *gabs.Container.
					ok = IsIncluded(search, element)
				case nil:
					// Element is not found in resulted JSON.
					ok = false
				default:
					if search.Data() == element.Data() {
						ok = true
					} else {
						ok = false
					}
					v = v
				}
				if ok == false {
					break
				}
			}
			// If underlying value is not an object.
		} else {
			if expected.Data() == resulted.Data() {
				ok = true
			} else {
				ok = false
			}
		}
	}
	return ok
}
