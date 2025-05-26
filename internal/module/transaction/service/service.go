package service

import (
	"context"
	"fmt"
	"math/big"
	"strconv"
	"strings"
	"time"

	"github.com/Digitalkeun-Creative/be-dzikra-pos-service/constants"
	"github.com/Digitalkeun-Creative/be-dzikra-pos-service/external/proto/transaction"
	"github.com/Digitalkeun-Creative/be-dzikra-pos-service/internal/module/transaction/dto"
	"github.com/Digitalkeun-Creative/be-dzikra-pos-service/internal/module/transaction/entity"
	transactionItemEntity "github.com/Digitalkeun-Creative/be-dzikra-pos-service/internal/module/transaction_item/entity"
	"github.com/Digitalkeun-Creative/be-dzikra-pos-service/pkg/err_msg"
	"github.com/Digitalkeun-Creative/be-dzikra-pos-service/pkg/utils"
	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog/log"
)

func (s *transactionService) CreateTransaction(ctx context.Context, req *dto.CreateTransactionRequest, tableNumber int) (*dto.CreateTransactionResponse, error) {
	memberResult, err := s.memberRepository.FindMemberByEmailAndPhoneNumber(ctx, req.Email, req.PhoneNumber)
	if err != nil {
		if strings.Contains(err.Error(), constants.ErrMemberNotFound) {
			log.Error().Err(err).Any("payload", req).Msg("service::CreateTransaction - Member not found")
			return nil, err_msg.NewCustomErrors(fiber.StatusNotFound, err_msg.WithMessage(constants.ErrMemberNotFound))
		}

		log.Error().Err(err).Msg("service::CreateTransaction - Failed to find member")
		return nil, err_msg.NewCustomErrors(fiber.StatusInternalServerError, err_msg.WithMessage(constants.ErrInternalServerError))
	}

	var isMember bool
	if memberResult != nil {
		isMember = true
	}

	status := constants.OrderStatusUnpaid
	if req.PaymentType == constants.OrderPaymentCash {
		status = constants.OrderStatusProcess
	}

	totalQty := big.NewInt(0)
	totalProdAmt := big.NewInt(0)
	totalProdCap := big.NewInt(0)

	var stockInputs []dto.MinIngredientInput

	for i := range req.TransactionItems {
		it := &req.TransactionItems[i]

		productResult, err := s.productRepository.FindProductByID(ctx, it.ProductID)
		if err != nil {
			if strings.Contains(err.Error(), constants.ErrProductNotFound) {
				log.Error().Err(err).Any("payload", req).Msg("service::CreateTransaction - Product not found")
				return nil, err_msg.NewCustomErrors(fiber.StatusNotFound, err_msg.WithMessage(constants.ErrProductNotFound))
			}

			log.Error().Err(err).Msg("service::CreateTransaction - Failed to find product")
			return nil, err_msg.NewCustomErrors(fiber.StatusInternalServerError, err_msg.WithMessage(constants.ErrInternalServerError))
		}

		qtyBI := big.NewInt(int64(it.Quantity))
		priceBI := big.NewInt(int64(productResult.RealPrice))
		capBI := big.NewInt(int64(productResult.CapitalPrice))
		subtotal := new(big.Int).Mul(priceBI, qtyBI)
		subCapTotal := new(big.Int).Mul(capBI, qtyBI)
		totalProdAmt.Add(totalProdAmt, subtotal)
		totalProdCap.Add(totalProdCap, subCapTotal)
		totalQty.Add(totalQty, qtyBI)

		stockInputs = append(stockInputs, dto.MinIngredientInput{
			ProductID: it.ProductID,
			Quantity:  it.Quantity,
		})

		it.ProductName = productResult.Name
		it.ProductPrice = int(subtotal.Uint64())
		it.ProductCapitalPrice = int(subCapTotal.Uint64())
		it.TotalAmount = int(subtotal.Uint64())
		it.TotalAmountCapitalPrice = int(subCapTotal.Uint64())
	}

	discountAmt := big.NewInt(0)

	if isMember {
		md, err := s.memberDiscountRepository.FindDiscount(ctx)
		if err != nil {
			log.Error().Err(err).Msg("service::CreateTransaction - Failed to find member discount")
			return nil, err_msg.NewCustomErrors(fiber.StatusInternalServerError, err_msg.WithMessage(constants.ErrInternalServerError))
		}

		pctBI := big.NewInt(int64(md.Discount))
		discountAmt = new(big.Int).Div(new(big.Int).Mul(totalProdAmt, pctBI), big.NewInt(100))
		req.DiscountPercentage = int(md.Discount)
	}

	taxCfg, err := s.taxRepository.FindTax(ctx)
	if err != nil {
		log.Error().Err(err).Msg("service::CreateTransaction - Failed to find tax")
		return nil, err_msg.NewCustomErrors(fiber.StatusInternalServerError, err_msg.WithMessage(constants.ErrInternalServerError))
	}

	taxAmt := new(big.Int).Div(new(big.Int).Mul(totalProdAmt, big.NewInt(int64(taxCfg.Tax))), big.NewInt(100))
	totalAfterDiscount := new(big.Int).Sub(totalProdAmt, discountAmt)
	grandTotal := new(big.Int).Add(totalAfterDiscount, taxAmt)

	transactionID, err := utils.GenerateUUIDv7String()
	if err != nil {
		log.Error().Err(err).Msg("service::CreateTransaction - Failed to generate transaction id")
		return nil, err_msg.NewCustomErrors(fiber.StatusInternalServerError, err_msg.WithMessage(constants.ErrInternalServerError))
	}

	createdAt := time.Now()
	var paymenType string

	if req.PaymentType == constants.OrderPaymentCash {
		paymenType = constants.OrderPaymentCash
	} else {
		paymenType = constants.OrderPaymentNonCash
	}

	payload := &entity.Transaction{
		ID:                       transactionID,
		Name:                     req.Name,
		Email:                    req.Email,
		PhoneNumber:              req.PhoneNumber,
		TableNumber:              tableNumber,
		IsMember:                 isMember,
		Status:                   status,
		PaymentType:              paymenType,
		TotalQuantity:            int(totalQty.Uint64()),
		TotalProductAmount:       int(totalProdAmt.Uint64()),
		TotalProductCapitalPrice: int(totalProdCap.Uint64()),
		DiscountPercentage:       req.DiscountPercentage,
		TaxAmount:                int(taxAmt.Uint64()),
		TotalAmount:              int(grandTotal.Uint64()),
		Notes:                    req.Notes,
		VPaymentID:               "",
		VPaymentRedirectUrl:      "",
		VTransactionID:           "",
		ChangeMoney:              0,
		TotalMoney:               req.TotalMoney,
		CreatedAt:                createdAt,
	}

	// Begin transaction
	tx, err := s.db.Beginx()
	if err != nil {
		log.Error().Err(err).Msg("service::CreateTransaction - Failed to begin transaction")
		return nil, err_msg.NewCustomErrors(fiber.StatusInternalServerError, err_msg.WithMessage(constants.ErrInternalServerError))
	}
	defer func() {
		if err != nil {
			if rollbackErr := tx.Rollback(); rollbackErr != nil {
				log.Error().Err(rollbackErr).Msg("service::CreateTransaction - Failed to rollback transaction")
			}
		}
	}()

	err = s.transactionRepository.InsertNewTransaction(ctx, tx, payload)
	if err != nil {
		log.Error().Err(err).Msg("service::CreateTransaction - insert new transaction failed")
		return nil, err_msg.NewCustomErrors(fiber.StatusInternalServerError, err_msg.WithMessage(constants.ErrInternalServerError))
	}

	var payloadItem *transactionItemEntity.TransactionItem
	for _, item := range req.TransactionItems {
		payloadItem = &transactionItemEntity.TransactionItem{
			Quantity:                item.Quantity,
			TotalAmount:             item.TotalAmount,
			ProductName:             item.ProductName,
			ProductPrice:            item.ProductPrice,
			TransactionID:           transactionID,
			ProductID:               item.ProductID,
			ProductCapitalPrice:     item.ProductCapitalPrice,
			TotalAmountCapitalPrice: item.TotalAmountCapitalPrice,
		}

		err := s.transactionItemRepository.InsertNewTransactionItem(ctx, tx, payloadItem)
		if err != nil {
			log.Error().Err(err).Msg("service::CreateTransaction - insert new transaction item failed")
			return nil, err_msg.NewCustomErrors(fiber.StatusInternalServerError, err_msg.WithMessage(constants.ErrInternalServerError))
		}
	}

	var res *transaction.CreateTransactionResponse
	if req.PaymentType == constants.OrderPaymentCash {
		paidStr, changeStr := utils.ComputeCash(req.TotalMoney, grandTotal)

		if err := s.transactionRepository.UpdateCashField(ctx, tx, payload.ID.String(), paidStr, changeStr); err != nil {
			log.Error().Err(err).Msg("service::CreateTransaction - update cash field failed")
			return nil, err_msg.NewCustomErrors(fiber.StatusInternalServerError, err_msg.WithMessage(constants.ErrInternalServerError))
		}

		go func() {
			if err = s.recipeService.MinIngredients(ctx, stockInputs); err != nil {
				log.Error().Err(err).Msg("service::CreateTransaction - failed to min ingredient stock")
			}
		}()

		// 	 go func() {
		//     notif := map[string]interface{}{
		//         "isUser": false,
		//         "title":  "Pesanan Baru Diterima",
		//         "detail": fmt.Sprintf("Pesanan %s baru dari meja %s sebesar Rp. %s", payload.ID, payload.TableNumber, paidStr),
		//         "url":    "",
		//     }
		//     b, _ := json.Marshal(notif)
		//     http.Post(os.Getenv("BASE_URL_ECOMERCE")+"/api/notification/send_fcm_batch", "application/json", bytes.NewBuffer(b))
		// }()
	} else {
		var transactionItems []*transaction.TransactionItem
		for _, item := range req.TransactionItems {
			transactionItems = append(transactionItems, &transaction.TransactionItem{
				Quantity:                fmt.Sprintf("%d", item.Quantity),
				TotalAmount:             fmt.Sprintf("%d", item.TotalAmount),
				ProductName:             item.ProductName,
				ProductPrice:            fmt.Sprintf("%d", item.ProductPrice),
				TransactionId:           transactionID.String(),
				ProductId:               int64(item.ProductID),
				ProductCapitalPrice:     fmt.Sprintf("%d", item.ProductCapitalPrice),
				TotalAmountCapitalPrice: fmt.Sprintf("%d", item.TotalAmountCapitalPrice),
			})
		}

		payloadTransaction := &transaction.CreateTransactionRequest{
			Name:               req.Name,
			Email:              req.Email,
			PhoneNumber:        req.PhoneNumber,
			CallbackFinish:     req.CallbackFinish,
			TransactionId:      transactionID.String(),
			TotalAmount:        int64(grandTotal.Uint64()),
			TransactionItems:   transactionItems,
			TaxAmount:          int64(taxAmt.Uint64()),
			TotalProductAmount: int64(totalProdAmt.Uint64()),
			DiscountPercentage: int64(req.DiscountPercentage),
		}

		res, err = s.externalTransaction.CreateTransaction(ctx, payloadTransaction)
		if err != nil {
			log.Error().Err(err).Msg("service::CreateTransaction - Failed to create transaction via external service")
			return nil, err_msg.NewCustomErrors(fiber.StatusInternalServerError, err_msg.WithMessage(constants.ErrInternalServerError))
		}

		if err = s.transactionRepository.UpdateTransactionByID(ctx, tx, &entity.Transaction{
			ID:                  transactionID,
			VTransactionID:      res.VTransactionId,
			VPaymentID:          res.VPaymentId,
			VPaymentRedirectUrl: res.VPaymentRedirectUrl,
		}); err != nil {
			log.Error().Err(err).Msg("service::CreateTransaction - Failed to update transaction")
			return nil, err_msg.NewCustomErrors(fiber.StatusInternalServerError, err_msg.WithMessage(constants.ErrInternalServerError))
		}
	}

	// Commit transaction
	if err := tx.Commit(); err != nil {
		log.Error().Err(err).Any("payload", req).Msg("service::CreateTransaction - Failed to commit transaction")
		return nil, err_msg.NewCustomErrors(fiber.StatusInternalServerError, err_msg.WithMessage(constants.ErrInternalServerError))
	}

	transactionResult, err := s.transactionRepository.FindTransactionWithItemsByID(ctx, transactionID.String())
	if err != nil {
		log.Error().Err(err).Msg("service::CreateTransaction - Failed to find transaction with items by ID")
		return nil, err_msg.NewCustomErrors(fiber.StatusInternalServerError, err_msg.WithMessage(constants.ErrInternalServerError))
	}

	var transactionItems []dto.TransactionItemResponse
	for _, item := range transactionResult.TransactionItems {
		transactionItems = append(transactionItems, dto.TransactionItemResponse{
			ID:                      item.ID,
			ProductID:               item.ProductID,
			ProductName:             item.ProductName,
			Quantity:                fmt.Sprintf("%d", item.Quantity),
			TotalAmount:             fmt.Sprintf("%d", item.TotalAmount),
			ProductPrice:            fmt.Sprintf("%d", item.ProductPrice),
			ProductCapitalPrice:     fmt.Sprintf("%d", item.ProductCapitalPrice),
			TotalAmountCapitalPrice: fmt.Sprintf("%d", item.TotalAmountCapitalPrice),
		})
	}

	var totalMoney *string
	if req.TotalMoney != 0 {
		s := strconv.Itoa(req.TotalMoney)
		totalMoney = &s
	} else {
		totalMoney = nil
	}

	var changeMoney *string
	if req.ChangeMoney != 0 {
		c := strconv.Itoa(req.ChangeMoney)
		changeMoney = &c
	} else {
		changeMoney = nil
	}

	return &dto.CreateTransactionResponse{
		ID:                       transactionID.String(),
		Status:                   status,
		PhoneNumber:              req.PhoneNumber,
		Name:                     req.Name,
		Email:                    req.Email,
		IsMember:                 isMember,
		TotalQuantity:            fmt.Sprintf("%d", totalQty.Uint64()),
		TotalProductAmount:       fmt.Sprintf("%d", totalProdAmt.Uint64()),
		TotalProductCapitalPrice: fmt.Sprintf("%d", totalProdCap.Uint64()),
		TotalAmount:              fmt.Sprintf("%d", grandTotal.Uint64()),
		DiscountPercentage:       fmt.Sprintf("%d", req.DiscountPercentage),
		PaymentType:              paymenType,
		TotalMoney:               totalMoney,
		ChangeMoney:              changeMoney,
		TableNumber:              req.TableNumber,
		Notes:                    req.Notes,
		TaxAmount:                fmt.Sprintf("%d", taxAmt.Uint64()),
		VTransactionID:           res.VTransactionId,
		VPaymentID:               res.VPaymentId,
		VPaymentRedirectUrl:      res.VPaymentRedirectUrl,
		CreatedAt:                utils.FormatTime(createdAt),
		TransactionItems:         transactionItems,
	}, nil
}

func (s *transactionService) GetListTransaction(ctx context.Context, page, limit int, search string) (*dto.GetListTransactionResponse, error) {
	// calculate pagination
	currentPage, perPage, offset := utils.Paginate(page, limit)

	// get list products
	transactions, total, err := s.transactionRepository.FindListTransaction(ctx, perPage, offset, search)
	if err != nil {
		log.Error().Err(err).Msg("service::GetListProduct - error getting list transactions")
		return nil, err_msg.NewCustomErrors(fiber.StatusInternalServerError, err_msg.WithMessage(constants.ErrInternalServerError))
	}

	// check if transactions is nil
	if transactions == nil {
		transactions = []dto.GetListTransaction{}
	}

	// calculate total pages
	totalPages := utils.CalculateTotalPages(total, perPage)

	// create map response
	response := dto.GetListTransactionResponse{
		Transactions: transactions,
		TotalPages:   totalPages,
		CurrentPage:  currentPage,
		PageSize:     perPage,
		TotalData:    total,
	}

	// return response
	return &response, nil
}
