package repository

import (
	"github.com/praveennagaraj97/shoppers-gocommerce/db"
	"github.com/praveennagaraj97/shoppers-gocommerce/pkg/color"
	logger "github.com/praveennagaraj97/shoppers-gocommerce/pkg/log"
	assetrepository "github.com/praveennagaraj97/shoppers-gocommerce/repository/asset"
	categoriesrepository "github.com/praveennagaraj97/shoppers-gocommerce/repository/categories"
	userrepository "github.com/praveennagaraj97/shoppers-gocommerce/repository/user"
	useraddressrepository "github.com/praveennagaraj97/shoppers-gocommerce/repository/useraddress"
	"go.mongodb.org/mongo-driver/mongo"
)

type Repositories struct {
	categoriesRepo *categoriesrepository.CategoriesRepository
	userRepo       *userrepository.UserRepository
	addressRepo    *useraddressrepository.UserAddressRepository
	assetsRepo     *assetrepository.AssetRepository
}

func (r *Repositories) Initialize(mongoClient *mongo.Client, dbName string) {
	// User Repo
	r.userRepo = &userrepository.UserRepository{}
	r.userRepo.InitializeRepository(db.OpenCollection(mongoClient, dbName, "user"))

	// address repo
	r.addressRepo = &useraddressrepository.UserAddressRepository{}
	r.addressRepo.InitializeRepository(db.OpenCollection(mongoClient, dbName, "userAddress"))

	// Categories repository
	r.categoriesRepo = &categoriesrepository.CategoriesRepository{}
	r.categoriesRepo.Initialize(db.OpenCollection(mongoClient, dbName, "categories"))

	// assets repo
	r.assetsRepo = &assetrepository.AssetRepository{}
	r.assetsRepo.Initialize(db.OpenCollection(mongoClient, dbName, "assets"))

	logger.PrintLog("Repositories initialized 📜", color.Gray)
}

func (r *Repositories) GetCategoriesRepo() *categoriesrepository.CategoriesRepository {
	return r.categoriesRepo
}

func (r *Repositories) GetUserRepo() *userrepository.UserRepository {
	return r.userRepo
}

func (r *Repositories) GetAddressRepo() *useraddressrepository.UserAddressRepository {
	return r.addressRepo
}

func (r *Repositories) GetAssetRepo() *assetrepository.AssetRepository {
	return r.assetsRepo
}
