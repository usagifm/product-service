package sqlx

import (
	"context"
	"encoding/json"
	"fmt"
	"math"
	"strconv"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/redis/go-redis/v9"
	"github.com/usagifm/product-service/config"
	"github.com/usagifm/product-service/model"
	"github.com/usagifm/product-service/utils/helper"
	"github.com/usagifm/product-service/utils/logger"
)

const (
	getProductById = iota
	countProduct
	createProduct
)

var (
	productQueries = []string{
		countProduct:   `SELECT COUNT(*) FROM products WHERE name LIKE $1;`,
		getProductById: `SELECT * FROM products WHERE id = $1;`,
		createProduct:  `INSERT INTO products (name,price,description,quantity) VALUES ($1,$2,$3,$4);`,
	}
)

type ProductRepository struct {
	db          *sqlx.DB
	stmts       []*sqlx.Stmt
	redisClient *redis.Client
	redisConfig *config.Redis
}

func InitProductRepositoryRepository(db *sqlx.DB, redisClient *redis.Client, redisConfig *config.Redis) (ProductRepository, error) {
	stmts, err := PrepareQueries(db, productQueries)
	if err != nil {
		return ProductRepository{}, err
	}
	return ProductRepository{
		stmts:       stmts,
		db:          db,
		redisClient: redisClient,
		redisConfig: redisConfig,
	}, nil
}

func (r ProductRepository) GetProductList(ctx context.Context, searchByName string, sortBy string, sortType string, limit int, offset int) (model.ProductList, error) {

	redisPath := helper.GetRedisDocumentName("product-list:search"+searchByName+":sort-by:"+sortBy+":sort-type:"+sortType+":limit:"+strconv.Itoa(limit)+":offset:"+strconv.Itoa(offset), r.redisConfig)

	productListRedis := r.redisClient.Get(ctx, redisPath).Val()
	if productListRedis != "" {
		logger.GetLogger(ctx).Infof("Get product list data from redis")
		var productList model.ProductList
		json.Unmarshal([]byte(productListRedis), &productList)
		return productList, nil
	}

	if sortType == "ASC" {
		sortType = "ASC"
	} else {
		sortType = "DESC"
	}

	var orderByClause string
	switch sortBy {
	case "created_at":
		orderByClause = "created_at " + sortType
	case "price":
		orderByClause = "price " + sortType
	case "name":
		orderByClause = "name " + sortType
	default:
		orderByClause = "created_at DESC" // Default sorting by newest
	}

	totalCount := 0
	errCount := r.stmts[countProduct].Get(&totalCount, "%"+searchByName+"%")
	if errCount != nil {
		return model.ProductList{}, errCount
	}

	// Custom query for easy customization
	query := `SELECT * FROM products WHERE name LIKE $1  ORDER BY ` + orderByClause + ` LIMIT $2 OFFSET $3`

	fmt.Println(query)

	products := []model.Product{}
	offsetQuery := 0 + (offset-1)*limit
	err := r.db.Select(&products, query, "%"+searchByName+"%", limit, offsetQuery)
	if err != nil {
		logger.GetLogger(ctx).Errorf("Get product list err:%v\n", err)
		return model.ProductList{}, err
	}

	numPages := int(math.Ceil(float64(totalCount) / float64(limit)))

	productList := model.ProductList{
		Limit:      limit,
		Offset:     offset,
		TotalPages: numPages,
		Products:   products,
	}

	if len(products) > 0 {
		helper.SetObjectToRedis(ctx, r.redisClient, productList, "Product List", redisPath, 5*time.Minute)
	}
	return productList, nil
}

func (r ProductRepository) GetProductById(ctx context.Context, productId int) (model.Product, error) {

	// Note : Do not do this in production, use better implementation with Elasticsearch or Meilisearch instead
	redisPath := helper.GetRedisDocumentName("product:"+strconv.Itoa(productId), r.redisConfig)

	productRedis := r.redisClient.Get(ctx, redisPath).Val()
	if productRedis != "" {
		logger.GetLogger(ctx).Infof("Get product data by id from redis")
		var product model.Product
		json.Unmarshal([]byte(productRedis), &product)
		return product, nil
	}

	product := model.Product{}
	err := r.stmts[getProductById].GetContext(ctx, &product, productId)
	logger.GetLogger(ctx).Infof("productId : %d", productId)
	if err != nil {
		logger.GetLogger(ctx).Errorf("Get product by id err:%v\n", err)
		return product, err
	}

	helper.SetObjectToRedis(ctx, r.redisClient, product, "Product Detail", redisPath, 5*time.Minute)

	return product, nil
}

func (r ProductRepository) CreateProduct(ctx context.Context, product model.Product) error {

	_, err := r.stmts[createProduct].ExecContext(ctx, product.Name, product.Price, product.Description, product.Quantity)
	if err != nil {
		logger.GetLogger(ctx).Errorf("Create product err:%v\n", err)
		return err
	}

	helper.DeleteAllRedisDocumentByPath(ctx, r.redisClient, r.redisConfig, "product-list")

	return nil
}
