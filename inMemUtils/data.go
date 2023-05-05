package inMemUtils

import (
	"github.com/Jeffail/gabs/v2"
	cmap "github.com/orcaman/concurrent-map"
	"strconv"
	"strings"
)

type StorageRecord struct {
	WalletId string `json:"wallet_id"`
	Id       string `json:"item_id"`
	Name     string `json:"name,omitempty"`
	Type     string `json:"type,omitempty"`
	Value    []byte `json:"value,omitempty"`
}

type Metadata struct {
	WalletId string `json:"wallet_id"`
	Value    string `json:"value"`
}

type TagsEncrypted struct {
	WalletId string `json:"wallet_id"`
	Name     string `json:"name"`
	Value    string `json:"value"`
	RecordId string `json:"record_id"`
}

type TagsPlaintext struct {
	WalletId string `json:"wallet_id"`
	Name     string `json:"name"`
	Value    string `json:"value"`
	RecordId string `json:"record_id"`
}

func (e *InMemoryStorage) getStorageMaps(walletHandle int) (walletId string, recordMap cmap.ConcurrentMap, tagsMapP cmap.ConcurrentMap, tagsMapE cmap.ConcurrentMap, errCode int) {
	handleKey := strconv.Itoa(walletHandle)
	wId, ok := e.WalletHandles.Get(handleKey)
	if !ok {
		return "", nil, nil, nil, 200 //200: "WalletInvalidHandle: Caller passed invalid wallet handle"
	}

	walletId = wId.(string)

	tmp, ok := e.StoredRecords.Get(walletId)
	if !ok {
		return "", nil, nil, nil, 210 //210: "WalletStorageError: Storage error occurred during wallet operation"
	}
	recordMap = tmp.(cmap.ConcurrentMap)

	tmp, ok = e.StoredTagsP.Get(walletId)
	if !ok {
		return "", nil, nil, nil, 210 //210: "WalletStorageError: Storage error occurred during wallet operation"
	}
	tagsMapP = tmp.(cmap.ConcurrentMap)

	tmp, ok = e.StoredTagsE.Get(walletId)
	if !ok {
		return "", nil, nil, nil, 210 //210: "WalletStorageError: Storage error occurred during wallet operation"
	}
	tagsMapE = tmp.(cmap.ConcurrentMap)

	return walletId, recordMap, tagsMapP, tagsMapE, 0
}

func (e *InMemoryStorage) addRecord(walletId string, record StorageRecord) (id string, err int) {
	tmp, ok := e.StoredRecords.Get(walletId)
	if !ok {
		return "", 210
	}

	recordMap := tmp.(cmap.ConcurrentMap)

	okC := false
	for t := range recordMap.IterBuffered() {
		sRecord := t.Val.(StorageRecord)
		if sRecord.Name == record.Name {
			okC = true
		}
	}

	if !okC { // record doesn't exist, prepare id
		_, id = e.RecordCounter.Get()
		record.Id = id
		recordMap.Set(id, record)
	} else { // duplicate
		return "", 213
	}

	e.StoredRecords.Set(walletId, recordMap)

	return id, 0
}

func (e *InMemoryStorage) addTagsP(walletId string, recordId string, tag TagsPlaintext) bool {
	lTagsP := make(map[TagsPlaintext][]string)

	tmp, ok := e.StoredTagsP.Get(walletId)
	if !ok {
		return ok
	}
	tagsMap := tmp.(cmap.ConcurrentMap)

	if !tagsMap.IsEmpty() {
		tmp, ok = tagsMap.Get(walletId)
		if !ok {
			return ok
		}

		lTagsP = tmp.(map[TagsPlaintext][]string)
	}
	lTagsP[tag] = append(lTagsP[tag], recordId)

	tagsMap.Set(walletId, lTagsP)
	e.StoredTagsP.Set(walletId, tagsMap)

	return ok
}

func (e *InMemoryStorage) addTagsE(walletId string, recordId string, tag TagsEncrypted) bool {
	lTagsE := make(map[TagsEncrypted][]string)

	tmp, ok := e.StoredTagsE.Get(walletId)
	if !ok {
		return ok
	}

	tagsMap := tmp.(cmap.ConcurrentMap)

	if !tagsMap.IsEmpty() {
		tmp, ok = tagsMap.Get(walletId)
		if !ok {
			return ok
		}

		lTagsE = tmp.(map[TagsEncrypted][]string)
	}
	lTagsE[tag] = append(lTagsE[tag], recordId)

	tagsMap.Set(walletId, lTagsE)
	e.StoredTagsE.Set(walletId, tagsMap)

	return ok
}

func (e *InMemoryStorage) rIdFromTags(tag string, walletId string, recordId string, tagsMapP cmap.ConcurrentMap, tagsMapE cmap.ConcurrentMap) bool {

	switch tag {
	case "": // if a tag wasn't passed, remove recordId from all tags
		if !tagsMapP.IsEmpty() {
			tmp, ok := tagsMapP.Get(walletId)
			if !ok {
				return false
			}

			lTagsP := tmp.(map[TagsPlaintext][]string)

			for t, s := range lTagsP {
				for index := len(s) - 1; index >= 0; index-- {
					if s[index] == recordId {
						s = append(s[:index], s[index+1:]...)
					}
				}

				if len(s) == 0 {
					delete(lTagsP, t)
				}
			}

			tagsMapP.Set(walletId, lTagsP)
			e.StoredTagsP.Set(walletId, tagsMapP)
		}

		if !tagsMapE.IsEmpty() {
			tmp, ok := tagsMapE.Get(walletId)
			if !ok {
				return false
			}

			lTagsE := tmp.(map[TagsEncrypted][]string)

			for t, s := range lTagsE {
				for index := len(s) - 1; index >= 0; index-- {
					if s[index] == recordId {
						s = append(s[:index], s[index+1:]...)
					}
				}

				if len(s) == 0 {
					delete(lTagsE, t)
				}
			}

			tagsMapE.Set(walletId, lTagsE)
			e.StoredTagsE.Set(walletId, tagsMapE)
		}
	default: // if a tag is passed, remove recordId from it
		if tag[0:1] == "~" {
			tmp, ok := tagsMapP.Get(walletId)
			if !ok {
				return false
			}

			lTagsP := tmp.(map[TagsPlaintext][]string)

			var tagP TagsPlaintext
			var s []string

			for tagP, s = range lTagsP {
				if tagP.Name == tag {
					for index := len(s) - 1; index >= 0; index-- {
						if s[index] == recordId {
							s = append(s[:index], s[index+1:]...)
						}
					}
					break
				}

				if len(s) == 0 {
					delete(lTagsP, tagP)
				}
			}

			lTagsP[tagP] = s
			tagsMapP.Set(walletId, lTagsP)
			e.StoredTagsP.Set(walletId, tagsMapP)
		} else {
			tmp, ok := tagsMapE.Get(walletId)
			if !ok {
				return false
			}

			lTagsE := tmp.(map[TagsEncrypted][]string)

			var tagE TagsEncrypted
			var s []string

			for tagE, s = range lTagsE {
				if tagE.Name == tag {
					for index := len(s) - 1; index >= 0; index-- {
						if s[index] == recordId {
							s = append(s[:index], s[index+1:]...)
						}
					}
					break
				}

				if len(s) == 0 {
					delete(lTagsE, tagE)
				}
			}
			lTagsE[tagE] = s
			tagsMapE.Set(walletId, lTagsE)
			e.StoredTagsE.Set(walletId, tagsMapE)
		}
	}

	return true
}

func (e *InMemoryStorage) andSearchCase(walletId string, query []*gabs.Container, recordMap cmap.ConcurrentMap, tagListE map[TagsEncrypted][]string, tagListP map[TagsPlaintext][]string) ([]StorageRecord, int) {
	var sRecords []StorageRecord
	var inters []string

	var notRecords []StorageRecord
	not := false

	ids := make(map[int][]string)

	for i, value := range query {
		for key, child := range value.ChildrenMap() {
			switch key {

			case "$and":
				andRecords, err := e.andSearchCase(walletId, child.Children(), recordMap, tagListE, tagListP)
				if err != 0 {
					return nil, err
				}

				sRecords = append(sRecords, andRecords...)
			case "$or":
				orRecords, err := e.orSearchCase(walletId, child.Children(), recordMap, tagListE, tagListP)
				if err != 0 {
					return nil, err
				}

				sRecords = append(sRecords, orRecords...)
			case "$not":
				records, err := e.notSearchCase(walletId, value.ChildrenMap(), recordMap, tagListE, tagListP)
				if err != 0 {
					return nil, err
				}

				notRecords = append(notRecords, records...)
				not = true
			default:
				sIds, ok := e.checkTags(walletId, value, tagListE, tagListP)
				if ok != 0 {
					return nil, ok
				}

				ids[i] = sIds
			}
		}
	}

	var s1, s2 []string
	i := 0

	for i = range ids {
		if inters != nil {
			s1 = inters
		} else {
			s1 = ids[i]
		}

		s2 = ids[i+1]
		if len(s2) != 0 {
			inters = intersect(s1, s2)
		}

		if len(inters) == 0 {
			return nil, 212
		}
	}

	for index := range inters {
		tmp, ok := recordMap.Get(inters[index])
		if !ok {
			return nil, 210 //210: "WalletStorageError: Storage error occurred during wallet operation"
		}

		record := tmp.(StorageRecord)
		sRecords = append(sRecords, record)
	}

	if not == true {
		for i = len(sRecords) - 1; i >= 0; i-- {
			for j := range notRecords {
				if len(sRecords) == 0 {
					break
				}

				if sRecords[i].Name == notRecords[j].Name && sRecords[i].Type == notRecords[j].Type {
					sRecords = append(sRecords[:i], sRecords[i+1:]...)
				}
			}
		}
	}
	return sRecords, 0
}

func (e *InMemoryStorage) orSearchCase(walletId string, query []*gabs.Container, recordMap cmap.ConcurrentMap, tagListE map[TagsEncrypted][]string, tagListP map[TagsPlaintext][]string) ([]StorageRecord, int) {
	var sRecords []StorageRecord
	var notRecords []StorageRecord

	var recordIds []string
	not := false

	for _, value := range query {
		for key, child := range value.ChildrenMap() {
			switch key {

			case "$and":
				andRecords, err := e.andSearchCase(walletId, child.Children(), recordMap, tagListE, tagListP)
				if err != 0 {
					return nil, err
				}

				sRecords = append(sRecords, andRecords...)

			case "$or":
				orRecords, err := e.orSearchCase(walletId, child.Children(), recordMap, tagListE, tagListP)
				if err != 0 {
					return nil, err
				}

				sRecords = append(sRecords, orRecords...)

			case "$not":
				records, err := e.notSearchCase(walletId, value.ChildrenMap(), recordMap, tagListE, tagListP)
				if err != 0 {
					return nil, err
				}

				notRecords = append(notRecords, records...)
				not = true
			default:
				sIds, ok := e.checkTags(walletId, value, tagListE, tagListP)
				if ok != 0 {
					return nil, ok
				}

				recordIds = append(recordIds, sIds...)
			}
		}
	}

	for i := range recordIds {
		tmp, ok := recordMap.Get(recordIds[i])
		if !ok {
			return nil, 210 //210: "WalletStorageError: Storage error occurred during wallet operation"
		}

		record := tmp.(StorageRecord)
		sRecords = append(sRecords, record)
	}

	if not == true {
		sRecords = append(sRecords, notRecords...)
	}

	return sRecords, 0
}

func (e *InMemoryStorage) notSearchCase(walletId string, query map[string]*gabs.Container, recordMap cmap.ConcurrentMap, tagListE map[TagsEncrypted][]string, tagListP map[TagsPlaintext][]string) ([]StorageRecord, int) {
	var sRecords []StorageRecord
	var inters []string

	for _, tag := range query {
		sIds, ok := e.checkTags(walletId, tag, tagListE, tagListP)
		if ok != 0 {
			return nil, ok
		}

		inters = append(inters, sIds...)
	}

	for i := range inters {
		tmp, ok := recordMap.Get(inters[i])
		if !ok {
			return nil, 210 //210: "WalletStorageError: Storage error occurred during wallet operation"
		}

		record := tmp.(StorageRecord)
		sRecords = append(sRecords, record)
	}

	return sRecords, 0
}

func (e *InMemoryStorage) checkTags(walletId string, query *gabs.Container, tagsMapE map[TagsEncrypted][]string, tagsMapP map[TagsPlaintext][]string) ([]string, int) {
	var recordIds []string

	for key, child := range query.ChildrenMap() {
		switch child.Data().(type) {
		case []interface{}:
			var ids []string

			for _, v := range child.Children() {
				sIds, ok := e.checkTags(walletId, v, tagsMapE, tagsMapP)
				if ok != 0 {
					return nil, ok
				}

				ids = append(ids, sIds...)
			}
			return ids, 0
		default:
			children := child.ChildrenMap() // check if it's an operator

			if len(children) != 0 {
				for op, tagVal := range child.ChildrenMap() {
					records, err := e.searchTag(op, key, tagVal, tagsMapE, tagsMapP)
					if err != 0 {
						return nil, err
					}

					return records, 0
				}
			} else {
				if key[0:1] == "~" { // plaintext
					for i, v := range tagsMapP {
						if i.Name == key && i.Value == child.Data().(string) {
							recordIds = append(recordIds, v...)
						}
					}
				} else { // encrypted
					for i, v := range tagsMapE {
						if i.Name == key && i.Value == child.Data().(string) {
							recordIds = append(recordIds, v...)
						}
					}
				}
			}
		}
	}

	return recordIds, 0
}

func (e *InMemoryStorage) searchTag(op string, tName string, tagValue *gabs.Container, tagsMapE map[TagsEncrypted][]string, tagsMapP map[TagsPlaintext][]string) ([]string, int) {
	var records []string

	tValue := tagValue.Data().(string)

	switch op {
	case "$neq":
		if tName[0:1] == "~" { // plaintext
			for i, v := range tagsMapP {

				if i.Name == tName && i.Value != tValue { //
					if c := intersect(records, v); c == nil {
						records = append(records, v...)
					}
				}
			}
		} else { // encrypted
			for i, v := range tagsMapE {
				if i.Name == tName && i.Value != tValue {
					if c := intersect(records, v); c == nil {
						records = append(records, v...)
					}
				}
			}
		}
	case "$gt":
		if tName[0:1] == "~" { // plaintext
			sTagV, errInt := strconv.Atoi(tValue)

			for i, v := range tagsMapP {
				if errInt == nil {
					tagV, _ := strconv.Atoi(i.Value)

					if i.Name == tName && tagV > sTagV {
						if c := intersect(records, v); c == nil {
							records = append(records, v...)
						}
					}
				} else {
					return nil, 214
				}
			}
		} else {
			return nil, 214
		}
	case "$gte":
		if tName[0:1] == "~" { // plaintext
			sTagV, err := strconv.Atoi(tValue)
			if err != nil {
				return nil, 214
			}

			for i, v := range tagsMapP {
				tagV, errV := strconv.Atoi(i.Value)
				if errV != nil {
					continue
				}

				if i.Name == tName && tagV >= sTagV { //
					if c := intersect(records, v); c == nil {
						records = append(records, v...)
					}
				}
			}
		} else {
			return nil, 214
		}
	case "$lt":
		if tName[0:1] == "~" { // plaintext
			sTagV, err := strconv.Atoi(tValue)
			if err != nil {
				return nil, 214
			}

			for i, v := range tagsMapP {
				tagV, errV := strconv.Atoi(i.Value)
				if errV != nil {
					return nil, 214
				}

				if i.Name == tName && tagV < sTagV { //
					if c := intersect(records, v); c == nil {
						records = append(records, v...)
					}
				}
			}
		} else {
			return nil, 214
		}
	case "$lte":
		if tName[0:1] == "~" { // plaintext
			sTagV, err := strconv.Atoi(tValue)
			if err != nil {
				return nil, 214
			}

			for i, v := range tagsMapP {
				tagV, errV := strconv.Atoi(i.Value)
				if errV != nil {
					return nil, 214
				}

				if i.Name == tName && tagV <= sTagV { //
					if c := intersect(records, v); c == nil {
						records = append(records, v...)
					}
				}
			}
		} else {
			return nil, 214
		}
	case "$like":
		if tName[0:1] == "~" { // plaintext

			if strings.Contains(tValue, "%") {
				strings.Replace(tValue, "%", "", 1)
			}

			for i, v := range tagsMapP {
				if i.Name == tName && strings.Contains(i.Value, tValue) {

					beforeCut, _, _ := strings.Cut(i.Value, tValue)
					if len(beforeCut) != 0 {
						return nil, 212
					}

					if c := intersect(records, v); c == nil {
						records = append(records, v...)
					}
				}
			}
		} else {
			return nil, 214
		}
	case "$in":
		var in []string
		for _, vals := range tagValue.Children() {
			in = append(in, vals.Data().(string))
		}

		if tName[0:1] == "~" { // plaintext
			for i, v := range tagsMapP {
				for _, inV := range in {
					if i.Name == tName && i.Value == inV {
						if c := intersect(records, v); c == nil {
							records = append(records, v...)
						}
					}
				}
			}
		} else { // encrypted
			for i, v := range tagsMapE {
				for _, inV := range in {
					if i.Name == tName && i.Value == inV {
						if c := intersect(records, v); c == nil {
							records = append(records, v...)
						}
					}
				}
			}
		}
	case "":
		if tName[0:1] == "~" { // plaintext
			for i, v := range tagsMapP {
				if i.Name == tName && i.Value == tValue {
					if c := intersect(records, v); c == nil {
						records = append(records, v...)
					}
				}
			}
		} else { // encrypted
			for i, v := range tagsMapE {
				if i.Name == tName && i.Value == tValue {
					if c := intersect(records, v); c == nil {
						records = append(records, v...)
					}
				}
			}
		}
	}

	if len(records) == 0 {
		return nil, 212
	} else {
		return records, 0
	}
}

func intersect(s1 []string, s2 []string) (c []string) {
	k := make(map[string]int)
	for _, num := range s1 {
		k[num]++
	}

	for _, num := range s2 {
		if k[num] > 0 {
			c = append(c, num)
			k[num]--
		}
	}

	return c
}
