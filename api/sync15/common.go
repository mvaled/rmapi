package sync15

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"os"
	"path"
	"sort"

	"github.com/juruen/rmapi/log"
)

const (
	cacheDirEnvVar = "RMAPI_CACHE_DIR"
)

func HashEntries(entries []*Entry) (string, error) {
	sort.Slice(entries, func(i, j int) bool { return entries[i].DocumentID < entries[j].DocumentID })
	hasher := sha256.New()
	for _, d := range entries {
		//TODO: back and forth converting
		bh, err := hex.DecodeString(d.Hash)
		if err != nil {
			return "", err
		}
		hasher.Write(bh)
	}
	hash := hasher.Sum(nil)
	hashStr := hex.EncodeToString(hash)
	return hashStr, nil
}

func getCachedTreePath() (string, error) {
	if cachedir := os.Getenv(cacheDirEnvVar); cachedir != "" {
		log.Trace.Println("Using cache directory RMAPI_CACHE_DIR", cachedir)
		rmapiFolder := path.Join(cachedir, "1.5")
		err := os.MkdirAll(rmapiFolder, 0700)
		if err != nil {
			return "", err
		}
		cacheFile := path.Join(rmapiFolder, "tree")
		return cacheFile, nil
	}

	cachedir, err := os.UserCacheDir()
	if err != nil {
		return "", err
	}
	rmapiFolder := path.Join(cachedir, "rmapi", "1.5")
	err = os.MkdirAll(rmapiFolder, 0700)
	if err != nil {
		return "", err
	}
	cacheFile := path.Join(rmapiFolder, "tree.cache")
	return cacheFile, nil
}

const cacheVersion = 3

func loadTree() (*HashTree, error) {
	cacheFile, err := getCachedTreePath()
	if err != nil {
		return nil, err
	}
	tree := &HashTree{}
	if _, err := os.Stat(cacheFile); err == nil {
		b, err := os.ReadFile(cacheFile)
		if err != nil {
			return nil, err
		}
		err = json.Unmarshal(b, tree)
		if err != nil {
			log.Error.Println("cache corrupt, resyncing")
			return tree, nil
		}
		if tree.CacheVersion != cacheVersion {
			log.Info.Println("wrong cache file version, resyncing")
			return &HashTree{}, nil
		}
	}
	log.Info.Println("cache loaded: ", cacheFile)

	return tree, nil
}

// save cached version of the tree
func saveTree(tree *HashTree) error {
	cacheFile, err := getCachedTreePath()
	log.Info.Println("Writing cache: ", cacheFile)
	if err != nil {
		return err
	}
	tree.CacheVersion = cacheVersion
	b, err := json.MarshalIndent(tree, "", "")
	if err != nil {
		return err
	}
	err = os.WriteFile(cacheFile, b, 0644)
	return err
}
