package service

import (
	"context"
	"fmt"
	"math"
	"strconv"
	"strings"

	"github.com/Digitalkeun-Creative/be-dzikra-pos-service/constants"
	ingredientEntity "github.com/Digitalkeun-Creative/be-dzikra-pos-service/internal/module/ingredient/entity"
	ingredientStockEntity "github.com/Digitalkeun-Creative/be-dzikra-pos-service/internal/module/ingredient_stock/entity"
	productEntity "github.com/Digitalkeun-Creative/be-dzikra-pos-service/internal/module/product/entity"
	"github.com/Digitalkeun-Creative/be-dzikra-pos-service/internal/module/recipe/dto"
	recipeEntity "github.com/Digitalkeun-Creative/be-dzikra-pos-service/internal/module/recipe/entity"
	transactionDto "github.com/Digitalkeun-Creative/be-dzikra-pos-service/internal/module/transaction/dto"
	"github.com/Digitalkeun-Creative/be-dzikra-pos-service/pkg/err_msg"
	"github.com/Digitalkeun-Creative/be-dzikra-pos-service/pkg/utils"
	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog/log"
)

func (s *recipeService) GetListRecipe(ctx context.Context, page, limit int, search string) (*dto.GetListRecipeResponse, error) {
	currentPage, perPage, offset := utils.Paginate(page, limit)

	products, total, err := s.productRepository.FindListProductRecipe(ctx, perPage, offset, search)
	if err != nil {
		log.Error().Err(err).Msg("service::GetListRecipe - error getting list recipe")
		return nil, err_msg.NewCustomErrors(fiber.StatusInternalServerError, err_msg.WithMessage(constants.ErrInternalServerError))
	}

	totalPages := utils.CalculateTotalPages(total, perPage)

	recIDs := []int{}
	for _, p := range products {
		if p.RecipeID != 0 {
			recIDs = append(recIDs, p.RecipeID)
		}
	}

	recs, err := s.recipeRepository.FindRecipeByIDs(ctx, recIDs)
	if err != nil {
		log.Error().Err(err).Msg("service::GetListRecipe - error getting recipe by ids")
		return nil, err_msg.NewCustomErrors(fiber.StatusInternalServerError, err_msg.WithMessage(constants.ErrInternalServerError))
	}

	ingredients, err := s.ingredientRepository.FindIngredientByRecipeIDs(ctx, recIDs)
	if err != nil {
		log.Error().Err(err).Msg("service::GetListRecipe - error find list ingredient by recipe ids")
		return nil, err_msg.NewCustomErrors(fiber.StatusInternalServerError, err_msg.WithMessage(constants.ErrInternalServerError))
	}

	stockIDs := make([]int, len(ingredients))
	for i, ing := range ingredients {
		stockIDs[i] = ing.IngredientStockID
	}

	ingredientStocks, err := s.ingredientStockRepository.FindIngredientStockByIDs(ctx, stockIDs)
	if err != nil {
		log.Error().Err(err).Msg("service::FindIngredientStockByIDs - error find list ingredient stock by stock ids")
		return nil, err_msg.NewCustomErrors(fiber.StatusInternalServerError, err_msg.WithMessage(constants.ErrInternalServerError))
	}

	stockMap := make(map[int]dto.IngredientStock, len(ingredientStocks))
	for _, st := range ingredientStocks {
		stockMap[st.ID] = dto.IngredientStock{
			ID:                       st.ID,
			Name:                     st.Name,
			RequiredStock:            st.RequiredStock,
			Unit:                     st.Unit,
			PricePerAmountStock:      st.PricePerAmountStock,
			AmountStockPerPrice:      st.AmountStockPerPrice,
			AmountPriceRequiredStock: st.AmountPriceRequiredStock,
		}
	}

	ingMap := make(map[int][]dto.Ingredient)
	for _, ing := range ingredients {
		ingMap[ing.RecipeID] = append(ingMap[ing.RecipeID], dto.Ingredient{
			ID:              ing.ID,
			Unit:            ing.Unit,
			Cost:            int(ing.Cost),
			RequiredStock:   int(ing.RequiredStock),
			RecipeID:        ing.RecipeID,
			IngredientStock: stockMap[ing.IngredientStockID],
		})
	}

	recMap := make(map[int]dto.Recipe, len(recs))
	for _, r := range recs {
		recMap[r.ID] = dto.Recipe{
			ID:           r.ID,
			CapitalPrice: r.CapitalPrice,
			Ingredients:  ingMap[r.ID],
		}
	}

	out := make([]dto.ProductWithRecipe, len(products))
	for i, p := range products {
		recipeDTO := dto.Recipe{
			ID:           0,
			CapitalPrice: 0,
			Ingredients:  []dto.Ingredient{},
		}
		if p.RecipeID > 0 {
			if r, ok := recMap[p.RecipeID]; ok {
				if r.Ingredients == nil {
					r.Ingredients = []dto.Ingredient{}
				}
				recipeDTO = r
			}
		}
		out[i] = dto.ProductWithRecipe{
			ID:      p.ID,
			Name:    p.Name,
			Recipes: recipeDTO,
		}
	}

	return &dto.GetListRecipeResponse{
		RecipesSerialize: out,
		TotalPages:       totalPages,
		CurrentPage:      currentPage,
		PageSize:         perPage,
		TotalData:        total,
	}, nil
}

func (s *recipeService) UpdateRecipe(ctx context.Context, req *dto.UpdateRecipeRequest, productID int) (*dto.UpdateRecipeResponse, error) {
	seen := map[int]bool{}

	product, err := s.productRepository.FindProductByID(ctx, productID)
	if err != nil {
		if strings.Contains(err.Error(), constants.ErrProductNotFound) {
			log.Error().Err(err).Any("payload", req).Msg("service::UpdateMember - Product not found")
			return nil, err_msg.NewCustomErrors(fiber.StatusNotFound, err_msg.WithMessage(constants.ErrProductNotFound))
		}

		log.Error().Err(err).Msg("service::UpdateRecipe - Failed to get product")
		return nil, err_msg.NewCustomErrors(fiber.StatusInternalServerError, err_msg.WithMessage(constants.ErrInternalServerError))
	}

	recipeID := product.RecipeID

	for _, ing := range req.Ingredients {
		ingredientID, err := strconv.Atoi(ing.IngredientID)
		if err != nil {
			log.Error().Msg("service::UpdateRecipe - error convert ingredient_id")
			return nil, err_msg.NewCustomErrors(fiber.StatusBadRequest, err_msg.WithMessage("error convert ingredient_id"))
		}

		if seen[ingredientID] {
			log.Error().Msg("service::UpdateRecipe - error duplicate ingredient_id in payload")
			return nil, err_msg.NewCustomErrors(fiber.StatusBadRequest, err_msg.WithMessage("duplicate ingredient_id in payload"))
		}
		seen[ingredientID] = true
	}

	for _, ing := range req.Ingredients {
		ingredientID, err := strconv.Atoi(ing.IngredientID)
		if err != nil {
			log.Error().Msg("service::UpdateRecipe - error convert ingredient_id")
			return nil, err_msg.NewCustomErrors(fiber.StatusBadRequest, err_msg.WithMessage("error convert ingredient_id"))
		}

		count, err := s.ingredientStockRepository.CountIngredientStockByID(ctx, ingredientID)
		if err != nil {
			log.Error().Err(err).Msg("service::UpdateRecipe - failed to count ingredient stock")
			return nil, err_msg.NewCustomErrors(fiber.StatusInternalServerError, err_msg.WithMessage(constants.ErrInternalServerError))
		}
		if count == 0 {
			log.Error().Err(err).Msg("service::UpdateRecipe - count ingredient stock is empty")
			return nil, err_msg.NewCustomErrors(fiber.StatusNotFound, err_msg.WithMessage(constants.ErrIngredientStockNotFound))
		}
	}

	// Begin transaction
	tx, err := s.db.Beginx()
	if err != nil {
		log.Error().Err(err).Msg("service::UpdateRecipe - Failed to begin transaction")
		return nil, err_msg.NewCustomErrors(fiber.StatusInternalServerError, err_msg.WithMessage(constants.ErrInternalServerError))
	}
	defer func() {
		if err != nil {
			if rollbackErr := tx.Rollback(); rollbackErr != nil {
				log.Error().Err(rollbackErr).Msg("service::UpdateRecipe - Failed to rollback transaction")
			}
		}
	}()

	if len(req.Ingredients) > 0 {
		if err = s.ingredientRepository.DeleteIngredientByRecipeID(ctx, tx, recipeID); err != nil {
			log.Error().Err(err).Msg("service::UpdateRecipe - Failed to delete ingredients")
			return nil, err_msg.NewCustomErrors(fiber.StatusInternalServerError, err_msg.WithMessage(constants.ErrInternalServerError))
		}
	}

	for _, ingredient := range req.Ingredients {
		ingredientID, err := strconv.Atoi(ingredient.IngredientID)
		if err != nil {
			log.Error().Msg("service::UpdateRecipe - error convert ingredient_id")
			return nil, err_msg.NewCustomErrors(fiber.StatusBadRequest, err_msg.WithMessage("error convert ingredient_id"))
		}

		if err = s.ingredientRepository.InsertNewIngredient(ctx, tx, &ingredientEntity.Ingredient{
			Unit:              ingredient.Unit,
			Cost:              float64(ingredient.Cost),
			RecipeID:          recipeID,
			IngredientStockID: ingredientID,
			RequiredStock:     float64(ingredient.RequiredStock),
		}); err != nil {
			log.Error().Err(err).Msg("service::UpdateRecipe - Failed to insert new ingredients")
			return nil, err_msg.NewCustomErrors(fiber.StatusInternalServerError, err_msg.WithMessage(constants.ErrInternalServerError))
		}
	}

	updatedIngs, err := s.ingredientRepository.FindIngredientByRecipeID(ctx, tx, recipeID)
	if err != nil {
		log.Error().Err(err).Msg("service::UpdateRecipe - Failed to get ingredient by recipe id")
		return nil, err_msg.NewCustomErrors(fiber.StatusInternalServerError, err_msg.WithMessage(constants.ErrInternalServerError))
	}

	var capitalPrice int
	for _, ui := range updatedIngs {
		capitalPrice += int(ui.Cost)
	}

	stockIDs := make([]int, len(updatedIngs))
	for i, ui := range updatedIngs {
		stockIDs[i] = ui.IngredientStockID
	}

	stocks, err := s.ingredientStockRepository.FindIngredientStockByIDs(ctx, stockIDs)
	if err != nil {
		log.Error().Err(err).Msg("service::UpdateRecipe - FindIngredientStocks")
		return nil, err_msg.NewCustomErrors(fiber.StatusInternalServerError, err_msg.WithMessage(constants.ErrInternalServerError))
	}

	stockMap := make(map[int]ingredientStockEntity.IngredientStock, len(stocks))
	for _, st := range stocks {
		stockMap[st.ID] = st
	}

	maxStock := math.MaxInt32
	for _, ui := range updatedIngs {
		avail := stockMap[ui.IngredientStockID].RequiredStock
		reqQty := int(ui.RequiredStock)
		if n := avail / reqQty; n < maxStock {
			maxStock = n
		}
	}

	if err = s.productRepository.UpdateProductStock(ctx, tx, &productEntity.Product{
		ID:    productID,
		Stock: maxStock,
	}); err != nil {
		log.Error().Err(err).Msg("service::UpdateRecipe - failed to update product stock")
		return nil, err_msg.NewCustomErrors(fiber.StatusInternalServerError, err_msg.WithMessage(constants.ErrInternalServerError))
	}

	if err = s.recipeRepository.UpdateRecipeCapitalPrice(ctx, tx, &recipeEntity.Recipe{
		ID:           recipeID,
		CapitalPrice: capitalPrice,
	}); err != nil {
		log.Error().Err(err).Msg("service::UpdateRecipe - failed to update recipe capital price")
		return nil, err_msg.NewCustomErrors(fiber.StatusInternalServerError, err_msg.WithMessage("internal error"))
	}

	if err = tx.Commit(); err != nil {
		log.Error().Err(err).Msg("service::UpdateRecipe - failed to commit transaction")
		return nil, err_msg.NewCustomErrors(fiber.StatusInternalServerError, err_msg.WithMessage(constants.ErrInternalServerError))
	}

	detail, err := s.recipeRepository.FindDetailRecipes(ctx, recipeID)
	if err != nil {
		log.Error().Err(err).Msg("service::UpdateRecipe - failed get detail recipes")
		return nil, err_msg.NewCustomErrors(fiber.StatusInternalServerError, err_msg.WithMessage(constants.ErrInternalServerError))
	}

	resp := &dto.UpdateRecipeResponse{
		ID:           detail.ID,
		CapitalPrice: detail.CapitalPrice,
		Product: dto.ProductDTO{
			ID:    detail.Product.ID,
			Name:  detail.Product.Name,
			Stock: detail.Product.Stock,
		},
		Ingredients: make([]dto.IngredientDetailDTO, len(detail.Ingredients)),
	}

	for i, ing := range detail.Ingredients {
		resp.Ingredients[i] = dto.IngredientDetailDTO{
			ID:            ing.ID,
			Unit:          ing.Unit,
			Cost:          ing.Cost,
			RequiredStock: fmt.Sprintf("%d", ing.RequiredStock),
			RecipeID:      ing.RecipeID,
			IngredientStock: dto.IngredientStockDTO{
				ID:   ing.Stock.ID,
				Name: ing.Stock.Name,
			},
		}
	}

	return resp, nil
}

func (s *recipeService) MinIngredients(ctx context.Context, inputs []transactionDto.MinIngredientInput) error {
	// 1. product → recipe
	productIDs := make([]string, len(inputs))
	for i, in := range inputs {
		productIDs[i] = strconv.Itoa(in.ProductID)
	}
	prods, err := s.productRepository.FindProductRecipeByProductIDs(ctx, productIDs)
	if err != nil {
		log.Error().Err(err).Msg("MinIngredients: FindProductRecipeByProductIDs failed")
		return err_msg.NewCustomErrors(fiber.StatusInternalServerError, err_msg.WithMessage(constants.ErrInternalServerError))
	}
	recMap := map[string]string{}
	for _, p := range prods {
		if p.RecipeID == 0 {
			return err_msg.NewCustomErrors(fiber.StatusBadRequest, err_msg.WithMessage("Some products do not have recipes"))
		}
		recMap[strconv.Itoa(p.ID)] = strconv.Itoa(p.RecipeID)
	}

	// 2. ambil bahan & stok (deduplikasi recipeIDs)
	seen := map[string]bool{}
	recipeIDs := make([]string, 0, len(recMap))
	for _, rid := range recMap {
		if !seen[rid] {
			recipeIDs = append(recipeIDs, rid)
			seen[rid] = true
		}
	}
	allIng, err := s.recipeRepository.FindIngredientsWithStock(ctx, recipeIDs)
	if err != nil {
		log.Error().Err(err).Msg("MinIngredients: FindIngredientsWithStock failed")
		return err_msg.NewCustomErrors(fiber.StatusInternalServerError, err_msg.WithMessage(constants.ErrInternalServerError))
	}

	// 3. decrement stok
	for _, in := range inputs {
		key := strconv.Itoa(in.ProductID)
		rid, ok := recMap[key]
		if !ok {
			return err_msg.NewCustomErrors(fiber.StatusBadRequest, err_msg.WithMessage("Invalid product ID"))
		}
		for _, ing := range allIng {
			if ing.RecipeID != rid {
				continue
			}
			need := ing.ReqPerUnit * in.Quantity
			if ing.StockAmount < need {
				return err_msg.NewCustomErrors(fiber.StatusBadRequest,
					err_msg.WithMessage(fmt.Sprintf("Stock insufficient for ingredient %s", ing.IngredientID)))
			}
			if derr := s.ingredientStockRepository.DecrementStock(ctx, ing.StockID, need); derr != nil {
				log.Error().Err(derr).Msg("MinIngredients: DecrementStock failed")
				return err_msg.NewCustomErrors(fiber.StatusInternalServerError, err_msg.WithMessage(constants.ErrInternalServerError))
			}
		}
	}

	// 4. ambil ulang stok terbaru
	allIng, err = s.recipeRepository.FindIngredientsWithStock(ctx, recipeIDs)
	if err != nil {
		log.Error().Err(err).Msg("MinIngredients: re-fetch stock failed")
		return err_msg.NewCustomErrors(fiber.StatusInternalServerError, err_msg.WithMessage(constants.ErrInternalServerError))
	}

	// 5. hitung maxStock per resep
	maxByRec := map[string]int{}
	for _, ing := range allIng {
		possible := ing.StockAmount / ing.ReqPerUnit
		old, ok := maxByRec[ing.RecipeID]
		if !ok || possible < old {
			maxByRec[ing.RecipeID] = possible
		}
	}

	// 6. map product→stock dan update
	prodStock := map[string]int{}
	for pid, rid := range recMap {
		prodStock[pid] = maxByRec[rid]
	}
	if err := s.productRepository.UpdateProductsStock(ctx, prodStock); err != nil {
		log.Error().Err(err).Msg("MinIngredients: UpdateProductsStock failed")
		return err_msg.NewCustomErrors(fiber.StatusInternalServerError, err_msg.WithMessage(constants.ErrInternalServerError))
	}

	return nil
}
