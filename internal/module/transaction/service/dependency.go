package service

import (
	externalNotification "github.com/Digitalkeun-Creative/be-dzikra-pos-service/external/notification"
	externalTransaction "github.com/Digitalkeun-Creative/be-dzikra-pos-service/external/transaction"
	externalUserFcmToken "github.com/Digitalkeun-Creative/be-dzikra-pos-service/external/user_fcm_token"
	memberPorts "github.com/Digitalkeun-Creative/be-dzikra-pos-service/internal/module/member/ports"
	memberDiscountPorts "github.com/Digitalkeun-Creative/be-dzikra-pos-service/internal/module/member_discount/ports"
	productPorts "github.com/Digitalkeun-Creative/be-dzikra-pos-service/internal/module/product/ports"
	recipePorts "github.com/Digitalkeun-Creative/be-dzikra-pos-service/internal/module/recipe/ports"
	taxPorts "github.com/Digitalkeun-Creative/be-dzikra-pos-service/internal/module/tax/ports"
	transactionPorts "github.com/Digitalkeun-Creative/be-dzikra-pos-service/internal/module/transaction/ports"
	transactionItemPorts "github.com/Digitalkeun-Creative/be-dzikra-pos-service/internal/module/transaction_item/ports"
	"github.com/jmoiron/sqlx"
)

var _ transactionPorts.TransactionService = &transactionService{}

type transactionService struct {
	db                        *sqlx.DB
	transactionRepository     transactionPorts.TransactionRepository
	transactionItemRepository transactionItemPorts.TransactionItemRepository
	memberRepository          memberPorts.MemberRepository
	productRepository         productPorts.ProductRepository
	memberDiscountRepository  memberDiscountPorts.MemberDiscountRepository
	taxRepository             taxPorts.TaxRepository
	recipeService             recipePorts.RecipeService
	externalTransaction       externalTransaction.ExternalTransaction
	externalNotification      externalNotification.ExternalNotification
	externalUserFcmToken      externalUserFcmToken.ExternalUserFcmToken
}

func NewTransactionService(
	db *sqlx.DB,
	transactionRepository transactionPorts.TransactionRepository,
	transactionItemRepository transactionItemPorts.TransactionItemRepository,
	memberRepository memberPorts.MemberRepository,
	productRepository productPorts.ProductRepository,
	memberDiscountRepository memberDiscountPorts.MemberDiscountRepository,
	taxRepository taxPorts.TaxRepository,
	recipeService recipePorts.RecipeService,
	externalTransaction externalTransaction.ExternalTransaction,
	externalNotification externalNotification.ExternalNotification,
	externalUserFcmToken externalUserFcmToken.ExternalUserFcmToken,
) *transactionService {
	return &transactionService{
		db:                        db,
		transactionRepository:     transactionRepository,
		transactionItemRepository: transactionItemRepository,
		memberRepository:          memberRepository,
		productRepository:         productRepository,
		memberDiscountRepository:  memberDiscountRepository,
		taxRepository:             taxRepository,
		recipeService:             recipeService,
		externalTransaction:       externalTransaction,
		externalNotification:      externalNotification,
		externalUserFcmToken:      externalUserFcmToken,
	}
}
