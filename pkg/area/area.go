package area

import (
	"encoding/json"
	"errors"
	"github.com/go-redis/redis/v8"
	"gorm.io/gorm"
	"log"
)

type RegionTool struct {
	db          *gorm.DB
	redisClient *redis.Client
	cityMap     map[string]*RegionInfo
	provinceMap map[string]*RegionInfo
}

type RegionOptions struct {
	DB          *gorm.DB
	RedisClient *redis.Client
	FilePath    string
}
type RegionInfo struct {
	RegionName   string `json:"regionName"`
	CityCode     string `json:"cityCode"`
	ProvinceCode string `json:"provinceCode"`
}

// DefaultRegionTool 初始化的时候有三种模式是否使用Redis，从文件中读取，从数据库中读取
// 1. Redis是否启用，指定redis的键
// 2. 初始化数据来源
func DefaultRegionTool(db *gorm.DB) *RegionTool {
	tool := &RegionTool{
		db: db,
	}
	go tool.loadFromDB()
	return tool
}

func NewRegionTool(op *RegionOptions) (tool *RegionTool, err error) {
	tool = &RegionTool{}
	if op.DB == nil {
		return nil, errors.New("DB is empty")
	}
	tool.db = op.DB
	if op.RedisClient != nil {
		// 去redis上查找键
		tool.redisClient = op.RedisClient
		result, err1 := tool.redisClient.Exists(nil, "region-city", "region-province").Result()
		if err1 != nil {
			return nil, err1
		}
		if result != 1 {
			go tool.loadFromDB()
		} else {
			tool.fullCityProvince()
		}
	}
	return
}

// loadFromDB 从数据库中加载数据, 并写入redis
func (tool *RegionTool) loadFromDB() {
	infos := make([]RegionInfo, 0)
	err := tool.db.Table("meta_region").Raw("").Find(&infos).Error
	if err != nil {
		log.Print(err)
		return
	}
	for _, info := range infos {
		tool.cityMap[info.CityCode] = &info
		if info.CityCode == "" {
			tool.provinceMap[info.ProvinceCode] = &info
		}
	}
	if tool.redisClient != nil {
		city, err := json.Marshal(tool.cityMap)
		if err != nil {
			log.Print(err)
			return
		}
		_, err = tool.redisClient.SetNX(nil, "region-city", city, 0).Result()
		if err != nil {
			log.Print(err)
			return
		}
		prov, err := json.Marshal(tool.provinceMap)
		if err != nil {
			log.Print(err)
			return
		}
		_, err = tool.redisClient.SetNX(nil, "region-province", prov, 0).Result()
		if err != nil {
			log.Print(err)
			return
		}
	}
	return
}

// fullCityProvince 从redis中加载数据
func (tool *RegionTool) fullCityProvince() {
	city, err := tool.redisClient.Get(nil, "region-city").Bytes()
	if err != nil {
		return
	}
	err = json.Unmarshal(city, &tool.cityMap)
	if err != nil {
		return
	}
	prov, err := tool.redisClient.Get(nil, "region-province").Bytes()
	if err != nil {
		return
	}
	err = json.Unmarshal(prov, &tool.cityMap)
	if err != nil {
		return
	}
}
